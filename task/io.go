package task

import (
	"encoding/json"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func isVersionFileExist() bool {
	if _, err := os.Stat(filepath.Join(config.Core.Workdir, "./version.json")); err != nil {
		return false
	}
	return true
}

func getCurrentVersionData() (result *utils.VersionData, err error) {
	var data []byte
	data, err = os.ReadFile(filepath.Join(config.Core.Workdir, "./version.json"))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, result)
	return
}

func getCurrentCategoriesList(version *utils.VersionData) (result *categoryUnitCollection, err error) {
	if version == nil {
		return nil, errors.New("version is nil")
	}
	var data []byte
	data, err = os.ReadFile(filepath.Join(config.Core.Workdir, version.Categories.Path))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, result)
	return
}

type bundleSentenceCollection = []bundleSentence

func getCurrentSentencesMap(categoriesList *categoryUnitCollection) (collection *map[string]bundleSentenceCollection, total int, err error) {
	if categoriesList == nil {
		return nil, 0, errors.New("categoriesList is nil")
	}
	collection = &map[string]bundleSentenceCollection{} // init
	for _, category := range *categoriesList {
		var data []byte
		data, err = os.ReadFile(filepath.Join(config.Core.Workdir, category.Path))
		if err != nil {
			return nil, 0, err
		}
		var result bundleSentenceCollection
		if err = json.Unmarshal(data, &result); err != nil {
			return nil, 0, err
		}
		(*collection)[category.Key] = result
		total += len(result)
	}
	return
}
