package task

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/utils"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

var (
	versionData    = &utils.VersionData{}
	categoriesData = &utils.CategoryUnitCollection{}
	sentencesMap   = &map[string]bundleSentenceCollection{}
	sentenceCount  = 0
)

// Start task loop
func Start() {
	logging.Logger.Info("Prepare Task environment...")
	if err := initTaskEnv(); err != nil {
		logging.Logger.Fatal(errors.WithMessage(err, "Can't init task env").Error())
	}
	// fmt.Printf("%+v", *categoriesData)
	logging.Logger.Info("Activate Loop: " + strconv.Itoa(config.Core.Interval) + "s at interval.")
	for {
		logging.Logger.Info("Sync remote repository...")
		if err := utils.SyncRepository(); err != nil {
			logging.Logger.Fatal(errors.WithMessage(err, "Can't sync remote repository").Error())
		}
		if err := Task(); err != nil {
			logging.Logger.Error(errors.WithMessage(err, "Exec task failure, will try in next loop").Error())
		}
		time.Sleep(time.Duration(config.Core.Interval) * time.Second)
	}
}

func initTaskEnv() error {
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
	logging.Logger.Info("Start Task. Initialize scoped env...")
	// 局部使用（深拷贝）
	scopedCurrentCategories := utils.CategoryUnitCollection{}.DeepCopy(*categoriesData)
	scopedCurrentSentences := categorySentenceBundleMap{}.DeepCopy(*sentencesMap)
	scopedCurrentVersionData := utils.VersionData{}.DeepCopy(*versionData)
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
		return errors.WithMessage(err, "Compare and update categories data failure.")
	} else if isCategoriesUpdated {
		scopedCurrentVersionData.UpdateCategoriesRecord(scopedCurrentCategories)
	} else {
		logging.Logger.Info("(2/2) Diff sentences data...")
		isSentencesUpdated := false
		for _, category := range *categoriesData {
			remoteSentences := (*remoteSentencesPointer)[category.Key]
			localSentences, ok := scopedCurrentSentences[category.Key]
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
