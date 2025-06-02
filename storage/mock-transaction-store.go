package storage

import (
	"fmt"
	"transaction-crud-svc-go-postgres/models"
)

type MockTransactionStore struct{}

var transactions map[uint]models.Transaction

func (store *MockTransactionStore) GetList() ([]models.Transaction, error) {
	var tx []models.Transaction

	for _, val := range transactions {
		tx = append(tx, val)
	}

	return tx, nil
}

func (store *MockTransactionStore) GetPaginatedList(limit int, offset int) ([]models.Transaction, error) {
	values := make([]models.Transaction, 0, len(transactions))
	for _, v := range transactions {
		values = append(values, v)
	}

	if offset > len(values) {
		return []models.Transaction{}, nil
	}

	end := offset + limit
	if end > len(values) {
		end = len(values)
	}

	return values[offset:end], nil
}

func (store *MockTransactionStore) Get(id uint) (models.Transaction, error) {
	val, exists := transactions[id]

	if exists {
		return val, nil
	} else {
		return val, fmt.Errorf("Not found")
	}
}

func (store *MockTransactionStore) Create(tx models.Transaction) error {
	transactions[tx.Id] = tx
	return nil
}

func (store *MockTransactionStore) Update(tx models.Transaction) error {
	transactions[tx.Id] = tx
	return nil
}

func (store *MockTransactionStore) Delete(id uint) error {
	delete(transactions, id)
	return nil
}

func (store *MockTransactionStore) DeleteAll() error {
	transactions = make(map[uint]models.Transaction)
	return nil
}
