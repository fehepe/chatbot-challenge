package mysql

import (
	"context"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB MySQL

type MySQL struct {
	DbClient *gorm.DB
	Ctx      context.Context
}

func ConnectDB(mysqlDSN string) error {
	context := context.Background()

	conn := strings.ReplaceAll(mysqlDSN, "chatdb", "")
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		return err
	}

	_ = db.Exec("CREATE DATABASE IF NOT EXISTS chatdb;")

	db, err = gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = MySQL{DbClient: db, Ctx: context}

	return nil

}
