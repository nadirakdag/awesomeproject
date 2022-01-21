package data

import (
	"errors"
)

type ActiveTab struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ActiveTabsRepository interface {
	GetAll() []ActiveTab
	Get(key string) *ActiveTab
	Add(item ActiveTab) error
}

var ErrActiveTabExists = errors.New("active tab key exists")

type ActiveTabsInMemoryRepository struct {
	ActiveTabs []ActiveTab
}

func NewActiveTabsInMemoryRepository() ActiveTabsRepository {
	return &ActiveTabsInMemoryRepository{
		ActiveTabs: []ActiveTab{},
	}
}

func (activeTab ActiveTabsInMemoryRepository) GetAll() []ActiveTab {
	return activeTab.ActiveTabs
}

func (activeTab ActiveTabsInMemoryRepository) Get(key string) *ActiveTab {

	for _, item := range activeTab.ActiveTabs {
		if item.Key == key {
			return &item
		}
	}

	return nil
}

func (activeTab *ActiveTabsInMemoryRepository) Add(newActiveTab ActiveTab) error {
	for _, item := range activeTab.ActiveTabs {
		if item.Key == newActiveTab.Key {
			return ErrActiveTabExists
		}
	}

	activeTab.ActiveTabs = append(activeTab.ActiveTabs, newActiveTab)
	return nil
}
