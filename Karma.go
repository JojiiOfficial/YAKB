package main

import (
	"fmt"

	"gorm.io/gorm"
)

// KarmaGiven represents given karma to a user
type KarmaGiven struct {
	db *DatabaseHandler `gorm:"-"`

	SenderID   int
	ReceiverID int
	MessageID  int
}

// UserKarma user account for karma
type UserKarma struct {
	db *DatabaseHandler `gorm:"-"`

	gorm.Model
	UserID   int
	UserName string
	Amount   int
}

// Save users karma
func (userKarma *UserKarma) Save() error {
	return userKarma.db.Save(userKarma).Error
}

// Save karmaGiven
func (karmaGiven *KarmaGiven) Save() error {
	return karmaGiven.db.Save(karmaGiven).Error
}

// Create karmaGiven
func (karmaGiven *KarmaGiven) Create() error {
	return karmaGiven.db.Create(karmaGiven).Error
}

func (userKarma UserKarma) String() string {
	return fmt.Sprintf("%s\t-\t%d\n", userKarma.UserName, userKarma.Amount)
}
