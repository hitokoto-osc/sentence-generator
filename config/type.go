package config

type core struct {
	Interval         int // Sync interval
	RemoteRepository string
	Workdir          string
}

type database struct {
	Driver             string
	Host               string
	Port               int
	User               string
	Password           string
	Database           string
	Charset            string
	CategoryTableName  string
	SentencesTableName string
}

type git struct {
	Name   string
	Email  string
	Branch string
	Driver string
	SSH    gitSSH
	HTTP   gitHTTP
}

type gitSSH struct {
	PrivateKey string
	Password   string
	User       string
}

type gitHTTP struct {
	User     string
	Password string
}

// Core config field
var Core *core

// Database config field
var Database *database

// Git config field
var Git *git
