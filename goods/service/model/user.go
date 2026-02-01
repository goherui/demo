package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(30);commentP:用户名"`
	Password string `gorm:"type:varchar(32);comment:密码"`
}

func (u *User) FindUser(db *gorm.DB, username string) error {
	return db.Where("username=?", username).First(&u).Error
}
