package task

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/utils"
)

func compareAndUpdateCategories(local, remote *utils.CategoryUnitCollection) (bool, error) {
	if len(*local) != len(*remote) {
		local = remote
		return true, nil
	}
	// 计算交集
	collection := utils.CategoryUnitCollection{}
	t := map[string]int{}
	for _, v := range *local {
		hash, err := getCategoryUnitHash(v)
		if err != nil {
			return false, err
		}
		t[hash] = 0
	}
	for _, v := range *remote {
		hash, err := getCategoryUnitHash(v)
		if err != nil {
			return false, err
		}
		if _, ok := t[hash]; ok {
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
	if len(*local) != len(*remote) {
		local = remote
		return true, nil
	}
	// 计算交集
	collection := bundleSentenceCollection{}
	t := map[string]int{}
	for _, v := range *local {
		hash, err := getBundleSentenceHash(v)
		if err != nil {
			return false, err
		}
		t[hash] = 0
	}
	for _, v := range *remote {
		hash, err := getBundleSentenceHash(v)
		if err != nil {
			return false, err
		}
		if _, ok := t[hash]; ok {
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
