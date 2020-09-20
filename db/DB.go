package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	// gorm.Model
	Challenge string `gorm:"primaryKey"`
	Id        string
}

type Publickey struct {
	Id        string `gorm:"primaryKey"`
	Publickey []byte
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
	if err := db.AutoMigrate(&Publickey{}).Error; err != nil {
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
	if err := db.Create(&User{Challenge: challenge, Id: id}).Error; err != nil {
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

func InsertPublicKey(id string, pubkeyJSON []byte) error {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return err
	}
	publickey := &Publickey{Id: id, Publickey: pubkeyJSON}
	if err := db.Create(publickey).Error; err != nil {
		var oldPubkey Publickey
		if err := db.Where("id = ?", id).Find(&oldPubkey).Update(&publickey).Error; err != nil {
			return err
		}
	}
	return nil
}
