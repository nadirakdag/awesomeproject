package handlers

import (
	"awesomeProject/data"
	"encoding/json"
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

	jE := json.NewEncoder(writer)
	jE.Encode(result)
}

func (inMemory *InMemory) getActiveTab(key string, writer http.ResponseWriter, request *http.Request) {

	activeTabResult := inMemory.activeTabsRepository.Get(key)

	jE := json.NewEncoder(writer)
	jE.Encode(activeTabResult)
}

func (inMemory *InMemory) addActiveTab(writer http.ResponseWriter, request *http.Request) {

	newActiveTab := &data.ActiveTab{}

	jD := json.NewDecoder(request.Body)
	err := jD.Decode(newActiveTab)

	if err != nil {
		http.Error(writer, "Payload format not valid", http.StatusBadRequest)
		return
	}

	err = inMemory.activeTabsRepository.Add(*newActiveTab)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	jE := json.NewEncoder(writer)
	jE.Encode(newActiveTab)
}
