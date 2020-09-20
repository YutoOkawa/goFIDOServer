package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	// gorm.Model
	Challenge string `gorm:"primaryKey"`
	UserID    string
}

func InitDB() error {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate((&User{})).Error; err != nil {
		return err
	}
	return nil
}

func InsertChallenge(challenge string, id string) error {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return err
	}
	if err := db.Create(&User{Challenge: challenge, UserID: id}).Error; err != nil {
		return err
	}
	return nil
}

func GetChallenge(challenge string) (User, error) {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return User{}, err
	}
	var user User
	if err := db.Where("challenge = ?", challenge).Find(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}
