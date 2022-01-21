package data

import (
	"awesomeProject/models"
	"errors"
)

type ActiveTabsRepository interface {
	GetAll() []models.ActiveTab
	Get(key string) *models.ActiveTab
	Add(item models.ActiveTab) error
}

var ErrActiveTabExists = errors.New("active tab key exists")

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

func (activeTab ActiveTabsInMemoryRepository) Get(key string) *models.ActiveTab {

	for _, item := range activeTab.ActiveTabs {
		if item.Key == key {
			return &item
		}
	}

	return nil
}

func (activeTab *ActiveTabsInMemoryRepository) Add(newActiveTab models.ActiveTab) error {
	for _, item := range activeTab.ActiveTabs {
		if item.Key == newActiveTab.Key {
			return ErrActiveTabExists
		}
	}

	activeTab.ActiveTabs = append(activeTab.ActiveTabs, newActiveTab)
	return nil
}
