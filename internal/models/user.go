package models

import (
	"github.com/fehepe/chatbot-challenge/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       uint   `json:"id" gorm:"unique"`
	Name     string `json:"name"`
	UserName string `json:"username" gorm:"unique"`
	Password []byte `json:"-"`
}

func NewUser(data map[string]string) (User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return User{}, err
	}

	user := User{
		Name:     data["name"],
		UserName: data["username"],
		Password: password,
	}

	mysql.DB.DbClient.Create(&user)

	return user, nil
}

func GetUserByUserName(data map[string]string) (User, error) {
	var user User

	mysql.DB.DbClient.Where("user_name = ?", data["username"]).First(&user)

	return user, nil
}

func GetUserById(id string) (User, error) {
	var user User

	mysql.DB.DbClient.Where("id = ?", id).First(&user)

	return user, nil
}

func AutoMigrate() {
	mysql.DB.DbClient.AutoMigrate(&User{})
}
