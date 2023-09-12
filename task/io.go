package task

import (
	"encoding/json"
	"github.com/cockroachdb/errors"
	"os"
	"path/filepath"

	"github.com/hitokoto-osc/sentence-generator/config"
	"github.com/hitokoto-osc/sentence-generator/database"
	"github.com/hitokoto-osc/sentence-generator/utils"
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
	result = &utils.VersionData{} // init
	err = json.Unmarshal(data, result)
	return
}

func getCurrentCategoriesList(version *utils.VersionData) (result *utils.CategoryUnitCollection, err error) {
	if version == nil {
		return nil, errors.WithStack(errors.New("version is nil"))
	}
	var data []byte
	data, err = os.ReadFile(filepath.Join(config.Core.Workdir, version.Categories.Path))
	if err != nil {
		return nil, err
	}
	result = &utils.CategoryUnitCollection{} // init
	err = json.Unmarshal(data, result)
	return
}

type bundleSentenceCollection []bundleSentence

// deepCopyBundleSentenceCollection deep copy from source
func deepCopyBundleSentenceCollection(collection bundleSentenceCollection) bundleSentenceCollection {
	tmp := bundleSentenceCollection{}
	for _, v := range collection {
		tmp = append(tmp, bundleSentence{
			Sentence: database.Sentence{
				ID:         v.ID,
				UUID:       v.UUID,
				Hitokoto:   v.Hitokoto,
				Type:       v.Type,
				From:       v.From,
				FromWho:    v.FromWho,
				Creator:    v.Creator,
				CreatorUID: v.CreatorUID,
				Reviewer:   v.Reviewer,
				CommitFrom: v.CommitFrom,
				CreatedAt:  v.CreatedAt,
			},
			Length: v.Length,
		})
	}
	return tmp
}

func getCurrentSentencesMap(categoriesList *utils.CategoryUnitCollection) (collection *map[string]bundleSentenceCollection, total int, err error) {
	if categoriesList == nil {
		return nil, 0, errors.WithStack(errors.New("categoriesList is nil"))
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
