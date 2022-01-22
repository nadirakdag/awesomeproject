package handlers

import (
	"awesomeProject/data"
	"awesomeProject/models"
	"log"
	"net/http"
)

type InMemory struct {
	l                    *log.Logger
	activeTabsRepository data.ActiveTabsRepository
}

func NewInMemory(l *log.Logger, activeTabsRepository data.ActiveTabsRepository) *InMemory {
	return &InMemory{l: l, activeTabsRepository: activeTabsRepository}
}

func (inMemory *InMemory) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {
		inMemory.l.Println("InMemory GET", request.URL)
		key := request.URL.Query().Get("key")
		if key == "" {
			inMemory.l.Println("InMemory GET, QueryParam Key is null returning all active tabs")
			inMemory.getActiveTabs(writer, request)
			return
		} else {
			inMemory.l.Println("InMemory GET, QueryParam Key is found returning active tab")
			inMemory.getActiveTab(key, writer, request)
			return
		}
	} else if request.Method == http.MethodPost {
		inMemory.addActiveTab(writer, request)
		return
	}

	writer.WriteHeader(http.StatusMethodNotAllowed)
}

func (inMemory *InMemory) getActiveTabs(writer http.ResponseWriter, request *http.Request) {

	result := inMemory.activeTabsRepository.GetAll()
	models.ActiveTabs(result).ToJSON(writer)
}

func (inMemory *InMemory) getActiveTab(key string, writer http.ResponseWriter, request *http.Request) {

	activeTabResult, err := inMemory.activeTabsRepository.Get(key)
	if err != nil {
		if err == data.ErrActiveTabNotFound {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = activeTabResult.ToJSON(writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (inMemory *InMemory) addActiveTab(writer http.ResponseWriter, request *http.Request) {

	newActiveTab := &models.ActiveTab{}

	if err := newActiveTab.FromJSON(request.Body); err != nil {
		http.Error(writer, "Payload format not valid", http.StatusBadRequest)
		return
	}

	if err := inMemory.activeTabsRepository.Add(*newActiveTab); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	err := newActiveTab.ToJSON(writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
