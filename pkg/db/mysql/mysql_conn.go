package mysql

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	dbClient *gorm.DB
	ctx      context.Context
}

func ConnectDB(mysqlDSN string) (*MySQL, error) {
	context := context.Background()

	mysqlClient, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &MySQL{dbClient: mysqlClient, ctx: context}, nil
}
