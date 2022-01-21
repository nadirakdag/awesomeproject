package data

import (
	"awesomeProject/models"
	"errors"
)

type ActiveTabsRepository interface {
	GetAll() []models.ActiveTab
	Get(key string) (*models.ActiveTab, error)
	Add(item models.ActiveTab) error
}

var ErrActiveTabExists = errors.New("active tab key exists")
var ErrActiveTabNotFound = errors.New("active tab does not exits")

type ActiveTabsInMemoryRepository struct {
	ActiveTabs []models.ActiveTab
}

func NewActiveTabsInMemoryRepository() ActiveTabsRepository {
	return &ActiveTabsInMemoryRepository{
		ActiveTabs: []models.ActiveTab{},
	}
}

func (activeTab ActiveTabsInMemoryRepository) GetAll() []models.ActiveTab {
	return activeTab.ActiveTabs
}

func (activeTab ActiveTabsInMemoryRepository) Get(key string) (*models.ActiveTab, error) {

	for _, item := range activeTab.ActiveTabs {
		if item.Key == key {
			return &item, nil
		}
	}

	return nil, ErrActiveTabNotFound
}

func (activeTab *ActiveTabsInMemoryRepository) Add(newActiveTab models.ActiveTab) error {
	tab, _ := activeTab.Get(newActiveTab.Key)
	if tab != nil {
		return ErrActiveTabExists
	}

	activeTab.ActiveTabs = append(activeTab.ActiveTabs, newActiveTab)
	return nil
}
