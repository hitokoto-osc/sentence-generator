package database

import (
	"fmt"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/upper/db/v4/adapter/mysql"
)

func mysqlConnect() (err error) {
	var mysqlSettings = mysql.ConnectionURL{
		User:     config.Database.User,
		Password: config.Database.Password,
		Database: config.Database.Database,
		Host:     fmt.Sprintf("%s:%v", config.Database.Host, config.Database.Port),
		Options: map[string]string{
			"charset": config.Database.Charset,
		},
	}

	Session, err = mysql.Open(mysqlSettings)
	if err != nil {
		return err
	}
	Session.SetMaxOpenConns(10)
	Session.SetMaxIdleConns(30)
	return nil
}
