package task

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"path/filepath"
	"unicode/utf8"
)

type categoryUnit struct {
	category
	Path string
}

type categoryUnitCollection []categoryUnit

// ImportFrom can convert []category to []categoryUnit
func (p *categoryUnitCollection) ImportFrom(c []category) {
	for _, v := range c {
		*p = append(*p, categoryUnit{
			category: v,
			Path:     filepath.Join(config.Core.Workdir, fmt.Sprintf("./sentences/%s.json", v.Key)),
		})
	}
}

func categorizeSentences(sentences []sentence) map[string]bundleSentenceCollection {
	result := map[string]bundleSentenceCollection{}
	for _, sentence := range sentences {
		if _, ok := result[sentence.Type]; !ok {
			result[sentence.Type] = bundleSentenceCollection{}
		}
		result[sentence.Type] = append(result[sentence.Type], bundleSentence{
			sentence: sentence,
			Length:   uint(utf8.RuneCountInString(sentence.Hitokoto)),
		})
	}
	return result
}

func getCategoryUnitHash(c categoryUnit) (string, error) {
	result, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	hash := MD5(result)
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
