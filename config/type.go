package config

type core struct {
	Interval         int    `mapstructure:"interval"` // Sync interval
	RemoteRepository string `mapstructure:"remote_repository"`
	Workdir          string `mapstructure:"workdir"`
}

type database struct {
	Driver             string `mapstructure:"driver"`
	Host               string `mapstructure:"host"`
	Port               int    `mapstructure:"port"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	Database           string `mapstructure:"database"`
	Charset            string `mapstructure:"charset"`
	CategoryTableName  string `mapstructure:"category_table_name"`
	SentencesTableName string `mapstructure:"sentences_table_name"`
}

type git struct {
	Name   string  `mapstructure:"name"`
	Email  string  `mapstructure:"email"`
	Branch string  `mapstructure:"branch"`
	Driver string  `mapstructure:"driver"`
	SSH    gitSSH  `mapstructure:"ssh"`
	HTTP   gitHTTP `mapstructure:"http"`
}

type gitSSH struct {
	PrivateKey string `mapstructure:"private_key"`
	Password   string `mapstructure:"password"`
	User       string `mapstructure:"user"`
}

type gitHTTP struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// Core config field
var Core *core

// Database config field
var Database *database

// Git config field
var Git *git
