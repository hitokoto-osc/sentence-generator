package task

import (
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"strconv"
	"time"

	"github.com/hitokoto-osc/sentence-generator/config"
	"github.com/hitokoto-osc/sentence-generator/logging"
	"github.com/hitokoto-osc/sentence-generator/utils"
)

var (
	versionData    = &utils.VersionData{}
	categoriesData = &utils.CategoryUnitCollection{}
	sentencesMap   = &map[string]bundleSentenceCollection{}
	sentenceCount  = 0
)

// Start task loop
func Start() {
	defer logging.Logger.Sync()
	logging.Logger.Info("Prepare Task environment...")
	if err := initTaskEnv(); err != nil {
		logging.Logger.Fatal("Can't init task env", zap.Error(err))
	}
	// fmt.Printf("%+v", *categoriesData)
	logging.Logger.Info("Activate Loop: " + strconv.Itoa(config.Core.Interval) + "s at interval.")
	for {
		logging.Logger.Info("Sync remote repository...")
		if err := utils.SyncRepository(); err != nil {
			logging.Logger.Fatal("Can't sync remote repository", zap.Error(err))
		}
		if err := Task(); err != nil {
			logging.Logger.Error("Exec task failure, will try in next loop", zap.Error(err))
		}
		time.Sleep(time.Duration(config.Core.Interval) * time.Second)
	}
}

func initTaskEnv() error {
	defer logging.Logger.Sync()
	var (
		isExist bool
		err     error
	)
	versionData, isExist, err = initVersionData()
	if err != nil {
		return err
	}
	if !isExist {
		logging.Logger.Info("[Task] 检测到版控数据不存在，使用默认数据初始化，如果您觉得这是错误的，请中断进程。任务将在 5 秒后自动开始...")
		time.Sleep(5 * time.Second)
	}
	categoriesData, err = initCategoriesData(isExist)
	if err != nil {
		return err
	}
	sentencesMap, sentenceCount, err = initSentencesMap(isExist)
	if err != nil {
		return err
	}
	sentenceCount = 0
	for _, sentence := range *sentencesMap {
		sentenceCount += len(sentence)
	}
	return nil
}

// Task is the main runner
func Task() error {
	defer logging.Logger.Sync()
	logging.Logger.Info("Start Task. Initialize scoped env...")
	// 局部使用（深拷贝）
	scopedCurrentCategories := utils.DeepCopyCategoryUnitCollection(*categoriesData)
	scopedCurrentSentences := deepCopyCategorySentenceBundleMap(*sentencesMap)
	scopedCurrentVersionData := utils.DeepCopyVersionData(*versionData)
	logging.Logger.Info("Fetch remote data...")
	// 获取远程数据
	// TODO: 合理利用 Total 提高比对效率
	remoteCategoriesPointer, remoteSentencesPointer, total, err := fetchRemoteData()
	if err != nil {
		return err
	}
	logging.Logger.Info("Start Diff process.")
	logging.Logger.Info("(1/2) Diff categories data...")
	isCategoriesUpdated, err := compareAndUpdateCategories(&scopedCurrentCategories, remoteCategoriesPointer)
	if err != nil {
		return errors.WithMessage(
			err,
			"Compare and update categories data failure",
		)
	} else if isCategoriesUpdated {
		scopedCurrentVersionData.UpdateCategoriesRecord(scopedCurrentCategories)
	} else {
		logging.Logger.Info("(2/2) Diff sentences data...")
		isSentencesUpdated := false
		for _, category := range *categoriesData {
			remoteSentences := (*remoteSentencesPointer)[category.Key]
			localSentences, ok := scopedCurrentSentences[category.Key]
			logging.Logger.Debug("Comparing category: " + category.Key)
			if !ok {
				scopedCurrentSentences[category.Key] = remoteSentences
				isSentencesUpdated = true
				continue
			}
			isUpdated, err := compareAndUpdateCategorySentences(&localSentences, &remoteSentences)
			if err != nil {
				return err
			} else if isUpdated {
				scopedCurrentSentences[category.Key] = remoteSentences
				if err = scopedCurrentVersionData.UpdateSentenceRecord(category.Key); err != nil {
					return err
				}
				isSentencesUpdated = true
			}
		}
		if !isSentencesUpdated {
			logging.Logger.Info("No change committed. Waiting for next loop.")
			return nil
		}
	}
	logging.Logger.Info("Changes checked. It will start generate process.")
	if err := scopedCurrentVersionData.StepVersion(); err != nil {
		return err
	}
	if err := doBundleRelease(
		&scopedCurrentCategories,
		&scopedCurrentSentences,
		&scopedCurrentVersionData,
	); err != nil {
		return err
	}
	logging.Logger.Info("Commit scoped variables to global cache...")
	categoriesData = &scopedCurrentCategories
	sentencesMap = &scopedCurrentSentences
	versionData = &scopedCurrentVersionData
	sentenceCount = total
	logging.Logger.Info("All jobs are finished. Waiting for next loop.")

	return nil
}
