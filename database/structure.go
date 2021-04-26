package database

import "time"

// Category is the database structure
type Category struct {
	ID        uint      `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Desc      string    `db:"desc" json:"desc"`
	Key       string    `db:"key" json:"key"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Sentence is the database structure
type Sentence struct {
	ID         uint    `db:"id" json:"id"  json:"id"`
	UUID       string  `db:"uuid"  json:"uuid"`
	Hitokoto   string  `db:"hitokoto"  json:"hitokoto"`
	Type       string  `db:"type"  json:"type"`
	From       string  `db:"from"  json:"from"`
	FromWho    *string `db:"from_who"  json:"from_who"`
	Creator    string  `db:"creator"  json:"creator"`
	CreatorUID uint    `db:"creator_uid"  json:"creator_uid"`
	Reviewer   uint    `db:"reviewer"  json:"reviewer"`
	CommitFrom string  `db:"commit_from"  json:"commit_from"`
	CreatedAt  string  `db:"created_at"  json:"created_at"`
}
