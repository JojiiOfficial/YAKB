package main

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DatabaseHandler contains data to store
type DatabaseHandler struct {
	*gorm.DB

	mx sync.Mutex
}

// NewDB creates a new db
func NewDB(file string) (*DatabaseHandler, error) {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automigration
	err = db.AutoMigrate(&KarmaGiven{}, &UserKarma{})
	if err != nil {
		return nil, err
	}

	return &DatabaseHandler{
		DB: db,
	}, nil
}

func (db *DatabaseHandler) addKarma(messageID, senderID int, receiver *tgbotapi.User, deltaAmount int) (bool, error) {
	db.mx.Lock()
	defer db.mx.Unlock()

	has, err := db.hasTransaction(messageID, senderID)
	if err != nil {
		return false, err
	}
	if has {
		fmt.Println("already gave karma to this message!")
		return false, nil
	}

	uk, err := db.getAccountByID(receiver)
	if err != nil {
		return false, err
	}

	uk.Amount = uk.Amount + deltaAmount

	err = db.addTransaction(senderID, receiver.ID, messageID)
	if err != nil {
		return false, err
	}

	return true, uk.Save()
}

func (db *DatabaseHandler) addTransaction(senderID, receiverID, messageID int) error {
	kg := KarmaGiven{
		db:         db,
		SenderID:   senderID,
		MessageID:  messageID,
		ReceiverID: receiverID,
	}

	return kg.Create()
}

// return true if karmatransaction exists
func (db *DatabaseHandler) hasTransaction(messageID, senderID int) (bool, error) {
	var i int64

	err := db.Model(&KarmaGiven{}).Where(&KarmaGiven{
		MessageID: messageID,
		SenderID:  senderID,
	}).Count(&i).Error

	if err != nil {
		return false, err
	}

	return i > 0, nil
}

func (db *DatabaseHandler) getAccountByID(user *tgbotapi.User) (*UserKarma, error) {
	var uk UserKarma

	hasAcc, err := db.hasAccount(user.ID)
	if err != nil {
		return nil, err
	}

	if hasAcc {
		// use existing acc
		err := db.Model(&UserKarma{}).Where(&UserKarma{
			UserID: user.ID,
		}).First(&uk).Error

		if err != nil {
			return nil, err
		}

		uk.db = db
		return &uk, nil
	}

	uk = UserKarma{
		db:       db,
		UserID:   user.ID,
		UserName: getNameFromUser(user),
	}

	// Create new acc
	err = db.Model(&UserKarma{}).Create(&uk).Error
	if err != nil {
		return nil, err
	}

	return &uk, nil
}

func (db *DatabaseHandler) hasAccount(userID int) (bool, error) {
	var i int64
	err := db.Model(&UserKarma{}).Where(&UserKarma{
		UserID: userID,
	}).Count(&i).Error

	if err != nil {
		return false, err
	}

	return i > 0, nil
}

func (db *DatabaseHandler) getTopKarma(userID int) ([]UserKarma, error) {
	var userkarmas []UserKarma
	c := *db.DB
	d := &c

	d = d.Model(&UserKarma{}).Limit(10).Order("amount desc")
	if userID > 0 {
		d = d.Where(&UserKarma{
			UserID: userID,
		})
	}

	err := d.Find(&userkarmas).Error
	if err != nil {
		return nil, err
	}

	return userkarmas, nil
}
