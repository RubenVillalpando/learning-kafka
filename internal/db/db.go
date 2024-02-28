package db

import (
	"fmt"
	"log"

	"github.com/RubenVillalpando/learning-kafka/internal/model"
	"github.com/google/uuid"
)

type DB struct {
	cache map[string]model.Customer
}

func New() *DB {
	return &DB{
		cache: make(map[string]model.Customer),
	}
}

func (db *DB) GetByID(id string) (*model.Customer, error) {
	if customer, ok := db.cache[id]; ok {
		log.Println(customer)
		return &customer, nil
	}
	return nil, fmt.Errorf("customer was not found with id: %s", id)
}

func (db *DB) StoreOne(customer *model.Customer) (string, error) {
	id := uuid.New().String()
	customer.ID = id
	db.cache[id] = *customer
	return id, nil
}
