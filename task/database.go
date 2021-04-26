package task

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/database"
	"time"
)

type category struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	Desc      string    `db:"desc"`
	Key       string    `db:"key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type sentence struct {
	ID         uint    `db:"id"`
	UUID       string  `db:"uuid"`
	Hitokoto   string  `db:"hitokoto"`
	Type       string  `db:"type"`
	From       string  `db:"from"`
	FromWho    *string `db:"from_who"`
	Creator    string  `db:"creator"`
	CreatorUID uint    `db:"creator_uid"`
	Reviewer   uint    `db:"reviewer"`
	CommitFrom string  `db:"commit_from"`
	CreatedAt  int64   `db:"created_at"`
}

type bundleSentence struct {
	sentence
	Length uint
}

// getCategories fetch categories from remote source set
func getCategories() ([]category, error) {
	col := database.
		Session.
		Collection(config.Database.CategoryTableName)
	if _, err := col.Exists(); err != nil { // check categoryTable whether is exist
		return nil, err
	}
	res := col.Find()
	var result []category
	err := res.All(&result)
	return result, err
}

// getCategories fetch sentences from remote source set
func getSentences() ([]sentence, error) {
	col := database.
		Session.
		Collection(config.Database.SentencesTableName)
	if _, err := col.Exists(); err != nil { // check sentencesTableName whether is exist
		return nil, err
	}
	res := col.Find()
	var result []sentence
	err := res.All(&result)
	return result, err
}
