package task

import (
	"encoding/json"
	"fmt"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/utils"
	"os"
	"path/filepath"
)

func initVersionData() (versionData *utils.VersionData, isExist bool, err error) {
	if !isVersionFileExist() {
		return &utils.VersionData{
			ProtocolVersion: "1.0.0",
			BundleVersion:   "1.0.0",
			UpdatedAt:       0,
			Categories: struct {
				Path      string `json:"path"`
				Timestamp int64  `json:"timestamp"`
			}{
				Path:      "./categories.json",
				Timestamp: 0,
			},
			Sentences: []utils.SentencesUnit{},
		}, false, nil
	}
	versionData, err = getCurrentVersionData()
	if err != nil {
		return nil, false, err
	}
	return versionData, true, nil
}

func initCategoriesData(isExist bool) (*utils.CategoryUnitCollection, error) {
	if !isExist {
		return &utils.CategoryUnitCollection{}, nil
	}
	return getCurrentCategoriesList(versionData)
}

func initSentencesMap(isExist bool) (*map[string]bundleSentenceCollection, int, error) {
	if !isExist {
		return nil, 0, nil
	}
	return getCurrentSentencesMap(categoriesData)
}

// TODO: optimize exports
//revive:disable:function-result-limit
func fetchRemoteData() (categories *utils.CategoryUnitCollection, sentencesMap *map[string]bundleSentenceCollection, total int, err error) {
	remoteCategoriesList, err := getCategories()
	if err != nil {
		return nil, nil, 0, err
	}
	categoriesList := utils.CategoryUnitCollection{}
	categoriesList.ImportFrom(remoteCategoriesList)
	remoteSentencesCollection, err := getSentences()
	if err != nil {
		return nil, nil, 0, err
	}
	total = len(remoteSentencesCollection)
	remoteSentencesMap := categorizeSentences(remoteSentencesCollection)
	return &categoriesList, &remoteSentencesMap, total, nil
}

//revive:enable:function-result-limit

func generateBundle(categories *utils.CategoryUnitCollection, sentencesMap *map[string]bundleSentenceCollection, versionData *utils.VersionData) error {
	logging.Logger.Info("Start generate sentences bundle.")
	logging.Logger.Info("(1/3) Write Version lockfile...")
	data, err := json.MarshalIndent(*versionData, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(config.Core.Workdir, "./version.json"), data, 666); err != nil {
		return err
	}
	logging.Logger.Info("(2/3) Write Categories data...")
	data, err = json.MarshalIndent(*categories, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(config.Core.Workdir, "./categories.json"), data, 666); err != nil {
		return err
	}
	logging.Logger.Info("(3/3) Write Sentences data...")
	for key, sentences := range *sentencesMap {
		logging.Logger.Info(fmt.Sprintf("Write category(%s) sentences...", key))
		data, err = json.MarshalIndent(sentences, "", "  ")
		if err != nil {
			return err
		}
		if err = os.WriteFile(
			filepath.Join(
				config.Core.Workdir,
				fmt.Sprintf("./sentences/%s.json", key),
			),
			data,
			666,
		); err != nil {
			return err
		}
	}
	logging.Logger.Info("Generate bundle successfully.")
	return nil
}

func doBundleRelease(categories *utils.CategoryUnitCollection, sentencesMap *map[string]bundleSentenceCollection, versionData *utils.VersionData) error {
	if err := generateBundle(categories, sentencesMap, versionData); err != nil {
		return err
	} // generate new build
	if err := utils.CommitAndPushRepository(); err != nil {
		return err
	}
	return utils.ReleaseAndPushRepository()
}
