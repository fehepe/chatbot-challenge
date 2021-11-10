package models

type User struct {
	Id       uint   `json:"id" gorm:"unique"`
	Name     string `json:"name"`
	UserName string `json:"username" gorm:"unique"`
	Password []byte `json:"-"`
}
