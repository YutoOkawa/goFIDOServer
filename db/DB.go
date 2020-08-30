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
	if err != nil {
		return err
	}
	db.AutoMigrate((&User{}))
	defer db.Close()
	return nil
}

func InsertDB(challenge string, id string) error {
	db, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	db.Create(&User{Challenge: challenge, UserID: id})
	defer db.Close()
	return nil
}

func DeleteDB(challenge string) error {
	db, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	var user User
	user.Challenge = challenge
	db.First(&user)
	db.Delete(&user)
	db.Close()
	return nil
}

func GetOneDB(challenge string) (User, error) {
	db, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		return User{}, err
	}
	var user User
	user.Challenge = challenge
	db.First(&user)
	db.Close()
	return user, nil
}
