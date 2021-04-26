package task

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/database"
)

type bundleSentence struct {
	database.Sentence
	Length uint `json:"length"`
}

// getCategories fetch categories from remote source set
func getCategories() ([]database.Category, error) {
	col := database.
		Session.
		Collection(config.Database.CategoryTableName)
	if _, err := col.Exists(); err != nil { // check categoryTable whether is exist
		return nil, err
	}
	res := col.Find()
	var result []database.Category
	err := res.All(&result)
	return result, err
}

// getCategories fetch sentences from remote source set
func getSentences() ([]database.Sentence, error) {
	col := database.
		Session.
		Collection(config.Database.SentencesTableName)
	if _, err := col.Exists(); err != nil { // check sentencesTableName whether is exist
		return nil, err
	}
	res := col.Find()
	var result []database.Sentence
	err := res.All(&result)
	return result, err
}
