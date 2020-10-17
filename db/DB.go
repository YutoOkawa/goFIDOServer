package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	gorm.Model
	Challenge string
	Userid    string
}

type Publickey struct {
	Keyid     string `gorm:"primaryKey"`
	Userid    string
	Username  string
	Publickey []byte
}

type Userdata struct {
	Userid    string `gorm:"primaryKey"`
	Username  string
	Signcount uint32
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
	if err := db.AutoMigrate(&Userdata{}).Error; err != nil {
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
	if err := db.Create(&User{Challenge: challenge, Userid: id}).Error; err != nil {
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

func DeleteChallenge(challenge string) error {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return err
	}
	user, err := GetChallenge(challenge)
	if err != nil {
		return err
	}
	if err := db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func InsertPublicKey(keyID string, userID string, userName string, pubkeyJSON []byte) error {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return err
	}
	publickey := &Publickey{Keyid: keyID, Userid: userID, Username: userName, Publickey: pubkeyJSON}
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
	if err := db.Where("username = ?", userID).Find(&pubkey).Error; err != nil {
		return Publickey{}, err
	}
	return pubkey, nil
}

func InsertUserData(userID string, userName string, signCount uint32) error {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return err
	}
	userData := &Userdata{Userid: userID, Username: userName, Signcount: signCount}
	if err := db.Create(userData).Error; err != nil {
		return err
	}
	return nil
}

func GetUserData(userID string) (Userdata, error) {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != err {
		return Userdata{}, err
	}

	var userdata Userdata
	if err := db.Where("userid = ?", userID).Find(&userdata).Error; err != nil {
		return Userdata{}, err
	}
	return userdata, nil
}

func UpdateSignCount(userID string, signCount uint32) error {
	db, err := gorm.Open("sqlite3", "users.db")
	defer db.Close()
	if err != nil {
		return err
	}

	if err := db.Model(&Userdata{}).Where("userid = ?", userID).Update("signcount", signCount).Error; err != nil {
		return err
	}
	return nil
}
