package task

import (
	"fmt"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/utils"
)

func compareAndUpdateCategories(local, remote *utils.CategoryUnitCollection) (bool, error) {
	defer logging.Logger.Sync()
	if len(*local) != len(*remote) {
		local = remote
		return true, nil
	}
	// 计算交集
	collection := utils.CategoryUnitCollection{}
	t := map[string]int{}
	for _, v := range *local {
		hash, err := getCategoryUnitHash(v)
		logging.Logger.Debug(fmt.Sprintf("Local exist: %s, hash: %s", v.Key, hash))
		if err != nil {
			return false, err
		}
		t[hash] = 0
	}
	for _, v := range *remote {
		hash, err := getCategoryUnitHash(v)
		logging.Logger.Debug(fmt.Sprintf("Remote exist: %s, hash: %s", v.Key, hash))
		if err != nil {
			return false, err
		}
		if _, ok := t[hash]; ok {
			logging.Logger.Debug(fmt.Sprintf("Category %s exist in both", v.Key))
			collection = append(collection, v)
		}
	}
	// 比对
	if len(collection) != len(*remote) {
		local = remote
		return true, nil
	}
	return false, nil
}

func compareAndUpdateCategorySentences(local, remote *bundleSentenceCollection) (bool, error) {
	defer logging.Logger.Sync()
	if len(*local) != len(*remote) {
		local = remote
		return true, nil
	}
	// 计算交集
	collection := bundleSentenceCollection{}
	t := map[string]int{}
	for _, v := range *local {
		hash, err := getBundleSentenceHash(v)
		logging.Logger.Debug(fmt.Sprintf("Local sentence: %+v \nhash: %s", v, hash))
		if err != nil {
			return false, err
		}
		t[hash] = 0
	}
	for _, v := range *remote {
		hash, err := getBundleSentenceHash(v)
		logging.Logger.Debug(fmt.Sprintf("Remote sentence: %+v \nhash: %s", v, hash))
		if err != nil {
			return false, err
		}
		if _, ok := t[hash]; ok {
			logging.Logger.Debug(fmt.Sprintf("Sentence %s exist in both", hash))
			collection = append(collection, v)
		}
	}
	// 比对
	if len(collection) != len(*remote) {
		local = remote
		return true, nil
	}
	return false, nil
}
