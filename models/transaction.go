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
	AccountId    string            `gorm:"index:idx_accountId;index:idx_accountId_createdAt" json:"accountId"`
	Type         TransactionType   `json:"type"`
	Amout        float64           `json:"amount"`
	Currency     string            `json:"currency"`
	CreatedAt    time.Time         `gorm:"index:idx_createdAt;index:idx_user_createdAt;index:idx_accountId_createdAt;index:idx_merchantId_createdAt" json:"createdAt"`
	ModifiedAt   time.Time         `json:"modifiedAt"`
	Description  string            `json:"description"`
	Status       TransactionStatus `gorm:"index:idx_status;index:idx_status_user;index:idx_status_accountId" json:"status"`
	MerchantId   string            `gorm:"index:idx_merchantId;index:idx_merchantId_createdAt" json:"merchantId"`
	MerchantName string            `json:"merchantName"`
	Metadata     string            `json:"metadata"`
	User         string            `gorm:"index:idx_user;index:idx_user_createdAt" json:"user"`
}

func MigrateTransaction(db *gorm.DB) error {

	err := db.AutoMigrate(&Transaction{})

	return err
}
