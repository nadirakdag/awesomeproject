package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type InMemory struct {
	l *log.Logger
}

type ActiveTabs struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewInMemory(l *log.Logger) *InMemory {
	return &InMemory{l: l}
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

	result := []*ActiveTabs{
		{Key: "active-tabs", Value: "Getir"},
		{Key: "info", Value: "Nadir Akdağ"},
	}

	jE := json.NewEncoder(writer)
	jE.Encode(result)
}

func (inMemory *InMemory) getActiveTab(key string, writer http.ResponseWriter, request *http.Request) {
	result := []*ActiveTabs{
		{Key: "active-tabs", Value: "Getir"},
		{Key: "info", Value: "Nadir Akdağ"},
	}

	activeTabResult := &ActiveTabs{}

	for _, activeTab := range result {
		if activeTab.Key == key {
			activeTabResult = activeTab
			break
		}
	}

	jE := json.NewEncoder(writer)
	jE.Encode(activeTabResult)
}

func (inMemory *InMemory) addActiveTab(writer http.ResponseWriter, request *http.Request) {

	newActiveTab := &ActiveTabs{}

	jD := json.NewDecoder(request.Body)
	err := jD.Decode(newActiveTab)

	if err != nil {
		http.Error(writer, "Payload format not valid", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(writer, "%v", newActiveTab)
}
