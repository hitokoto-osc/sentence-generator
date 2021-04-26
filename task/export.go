package task

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/utils"
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"
	"time"
)

var (
	versionData    *utils.VersionData
	categoriesData *categoryUnitCollection
	sentencesMap   *map[string]bundleSentenceCollection
	sentenceCount  = 0
)

// Start task loop
func Start() {
	for {
		logging.Logger.Info("Prepare Task environment...")
		if err := initTaskEnv(); err != nil {
			logging.Logger.Fatal(errors.WithMessage(err, "Can't init task env").Error())
		}
		if err := Task(); err != nil {
			logging.Logger.Error(errors.WithMessage(err, "Exec task failure, will try in next loop.").Error())
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
	scopedCurrentCategories := deepcopy.Copy(*categoriesData).(categoryUnitCollection)
	scopedCurrentSentences := deepcopy.Copy(*sentencesMap).(map[string]bundleSentenceCollection)
	scopedCurrentVersionData := deepcopy.Copy(*versionData).(utils.VersionData)
	logging.Logger.Info("Fetch remote data...")
	// 获取远程数据
	// TODO: 合理利用 Total 提高比对效率
	remoteCategoriesPointer, remoteSentencesPointer, total, err := fetchRemoteData()
	logging.Logger.Info("Start Diff process.")
	logging.Logger.Info("(1/2) Diff categories data...")
	if err != nil {
		return err
	}
	isCategoriesUpdated, err := compareAndUpdateCategories(&scopedCurrentCategories, remoteCategoriesPointer)
	if err != nil {
		return errors.WithMessage(err, "Compare and update categories data failure.")
	} else if isCategoriesUpdated {
		scopedCurrentVersionData.UpdateCategoriesRecord()
	}
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
	if isCategoriesUpdated || isSentencesUpdated {
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
	} else {
		logging.Logger.Info("No change committed. Waiting for next loop.")
	}
	return nil
}
