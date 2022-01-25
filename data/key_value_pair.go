package data

import (
	"awesomeProject/models"
	"errors"
)

// Key value pair repository interface
type KeyValueRepository interface {
	GetAll() models.KeyValuePairs
	Get(key string) (*models.KeyValuePair, error)
	Add(item models.KeyValuePair) error
}

var ErrKeyValuePairExists = errors.New("key value pair key exists")
var ErrKeyValuePairNotFound = errors.New("key value pair does not exits")

//in-memory database repository for Key value pairs
type KeyValueInMemoryRepository struct {
	KeyValues []models.KeyValuePair
}

// creates a new in-memory database repository for Key value pairs
func NewKeyValueInMemoryRepository() KeyValueRepository {
	return &KeyValueInMemoryRepository{
		KeyValues: []models.KeyValuePair{},
	}
}

// implements GetAll method from KeyValueRepository for in-memory database
// returns all Key value pairs from in-memory database
func (keyValue KeyValueInMemoryRepository) GetAll() models.KeyValuePairs {
	return keyValue.KeyValues
}

// implements Get method from KeyValueRepository for in-memory database
// returns an Key value pair that request with key from in-memory database
// if tab does not exists then returns ErrKeyValuePairNotFound
func (keyValue KeyValueInMemoryRepository) Get(key string) (*models.KeyValuePair, error) {

	for _, item := range keyValue.KeyValues {
		if item.Key == key {
			return &item, nil
		}
	}

	return nil, ErrKeyValuePairNotFound
}

// implements Add method from KeyValueRepository for in-memory database
// trys to add new Key value pair to in-memory database
// if tab already exists returns ErrKeyValuePairExists error
func (keyValue *KeyValueInMemoryRepository) Add(newKeyValuePair models.KeyValuePair) error {
	tab, _ := keyValue.Get(newKeyValuePair.Key)
	if tab != nil {
		return ErrKeyValuePairExists
	}

	keyValue.KeyValues = append(keyValue.KeyValues, newKeyValuePair)
	return nil
}
