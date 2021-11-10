package mysql

import (
	"context"

	"github.com/fehepe/chatbot-challenge/internal/models"
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

	mysqlClient, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = MySQL{DbClient: mysqlClient, Ctx: context}

	DB.DbClient.AutoMigrate(&models.User{})

	return nil

}
