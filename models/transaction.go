package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType int
type TransactionStatus int

const (
	UnknownType TransactionType = iota
	Credit
	Debit
)

const (
	UnknownStatus TransactionStatus = iota
	Completed
	Pending
	Failed
)

type TxMetadata struct {
	Channel  string `json:"channel"`
	Location string `json:"location"`
}

type Transaction struct {
	Id           uint              `gorm:"primary key;autoIncrement" json:"id"`
	AccountId    string            `json:"accountId"`
	Type         TransactionType   `json:"type"`
	Amout        float64           `json:"amount"`
	Currency     string            `json:"currency"`
	CreatedAt    time.Time         `json:"createdAt"`
	ModifiedAt   time.Time         `json:"modifiedAt"`
	Description  string            `json:"description"`
	Status       TransactionStatus `json:"status"`
	MerchantId   string            `json:"merchantId"`
	MerchantName string            `json:"merchantName"`
	Metadata     string            `json:"metadata"`
	User         string            `json:"user"`
}

func MigrateTransaction(db *gorm.DB) error {

	err := db.AutoMigrate(&Transaction{})

	return err
}
