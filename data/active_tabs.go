package data

import (
	"awesomeProject/models"
	"errors"
)

// active tab repository interface
type ActiveTabsRepository interface {
	GetAll() []models.ActiveTab
	Get(key string) (*models.ActiveTab, error)
	Add(item models.ActiveTab) error
}

var ErrActiveTabExists = errors.New("active tab key exists")
var ErrActiveTabNotFound = errors.New("active tab does not exits")

// active tabs in-memory database repository for active tabs
type ActiveTabsInMemoryRepository struct {
	ActiveTabs []models.ActiveTab
}

// creates a new in-memory database repository for active tabs
func NewActiveTabsInMemoryRepository() ActiveTabsRepository {
	return &ActiveTabsInMemoryRepository{
		ActiveTabs: []models.ActiveTab{},
	}
}

// implements GetAll method from ActiveTabsRepository for in-memory database
// returns all active tabs from in-memory database
func (activeTab ActiveTabsInMemoryRepository) GetAll() []models.ActiveTab {
	return activeTab.ActiveTabs
}

// implements Get method from ActiveTabsRepository for in-memory database
// returns an active tab that request with key from in-memory database
// if tab does not exists then returns ErrActiveTabNotFound
func (activeTab ActiveTabsInMemoryRepository) Get(key string) (*models.ActiveTab, error) {

	for _, item := range activeTab.ActiveTabs {
		if item.Key == key {
			return &item, nil
		}
	}

	return nil, ErrActiveTabNotFound
}

// implements Add method from ActiveTabsRepository for in-memory database
// trys to add new active tab to in-memory database
// if tab already exists returns ErrActiveTabExists error
func (activeTab *ActiveTabsInMemoryRepository) Add(newActiveTab models.ActiveTab) error {
	tab, _ := activeTab.Get(newActiveTab.Key)
	if tab != nil {
		return ErrActiveTabExists
	}

	activeTab.ActiveTabs = append(activeTab.ActiveTabs, newActiveTab)
	return nil
}
