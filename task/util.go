package task

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/hitokoto-osc/sentence-generator/database"
	"github.com/hitokoto-osc/sentence-generator/utils"
)

type categorySentenceBundleMap map[string]bundleSentenceCollection

// deepCopyCategorySentenceBundleMap deep copy from source
func deepCopyCategorySentenceBundleMap(collection categorySentenceBundleMap) map[string]bundleSentenceCollection {
	tmp := categorySentenceBundleMap{}
	for k, v := range collection {
		tmp[k] = deepCopyBundleSentenceCollection(v)
	}
	return tmp
}

func categorizeSentences(sentences []database.Sentence) map[string]bundleSentenceCollection {
	result := map[string]bundleSentenceCollection{}
	for _, sentence := range sentences {
		if _, ok := result[sentence.Type]; !ok {
			result[sentence.Type] = bundleSentenceCollection{}
		}
		result[sentence.Type] = append(result[sentence.Type], bundleSentence{
			Sentence: sentence,
			Length:   uint(utf8.RuneCountInString(sentence.Hitokoto)),
		})
	}
	return result
}

func getCategoryUnitHash(c utils.CategoryUnit) (string, error) {
	hash := MD5([]byte(fmt.Sprintf("%s.%s.%s.%s.%s", strconv.Itoa(int(c.ID)), c.Name, c.Path, c.Key, c.Desc)))
	return hex.EncodeToString(hash[:]), nil
}

func getBundleSentenceHash(c bundleSentence) (string, error) {
	result, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	hash := MD5(result)
	return hex.EncodeToString(hash[:]), nil
}

// MD5 calc md5
func MD5(s []byte) [16]byte {
	return md5.Sum(s)
}
