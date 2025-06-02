package storage

import (
	"transaction-crud-svc-go-postgres/models"

	"gorm.io/gorm"
)

type TransactionStore interface {
	GetList() ([]models.Transaction, error)
	Get(id uint) (models.Transaction, error)
	Create(tx models.Transaction) error
	Update(tx models.Transaction) error
	Delete(id uint) error
	DeleteAll() error
}

type PostgresTransactionStore struct {
	DB *gorm.DB
}

func (store *PostgresTransactionStore) GetList() ([]models.Transaction, error) {
	txModels := []models.Transaction{}

	err := store.DB.Find(&txModels).Error

	return txModels, err
}

func (store *PostgresTransactionStore) Get(id uint) (models.Transaction, error) {
	txModel := models.Transaction{}

	err := store.DB.Where("id = ?", id).First(&txModel).Error

	return txModel, err
}

func (store *PostgresTransactionStore) Create(tx models.Transaction) error {
	return store.DB.Create(&tx).Error
}

func (store *PostgresTransactionStore) Update(tx models.Transaction) error {
	return store.DB.Save(&tx).Error
}

func (store *PostgresTransactionStore) Delete(id uint) error {
	tx := &models.Transaction{}

	return store.DB.Delete(tx, id).Error
}

func (store *PostgresTransactionStore) DeleteAll() error {
	return store.DB.Exec("truncate table transactions").Error
}
