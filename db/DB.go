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
	Keyid     string `gorm:"primaryKey"`
	Userid    string
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

func InsertPublicKey(keyID string, userID string, pubkeyJSON []byte) error {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return err
	}
	publickey := &Publickey{Keyid: keyID, Userid: userID, Publickey: pubkeyJSON}
	if err := db.Create(publickey).Error; err != nil {
		return err
	}
	return nil
}

func GetPublicKey(userID string) (Publickey, error) {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return Publickey{}, err
	}

	var pubkey Publickey
	if err := db.Where("userid = ?", userID).Find(&pubkey).Error; err != nil {
		return Publickey{}, err
	}
	return pubkey, nil
}
