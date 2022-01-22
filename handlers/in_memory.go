package handlers

import (
	"awesomeProject/data"
	"awesomeProject/models"
	"log"
	"net/http"
)

// is a http handler object
type InMemory struct {
	l                    *log.Logger
	activeTabsRepository data.ActiveTabsRepository
}

// Creates a new InMemory handler given logger and active tab repository
func NewInMemory(l *log.Logger, activeTabsRepository data.ActiveTabsRepository) *InMemory {
	return &InMemory{l: l, activeTabsRepository: activeTabsRepository}
}

// @Summary ServeHTTP is the main entry point for the handler and staisfies the http.Handler
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

// @Summary Get All Active Tabs
// @Produce json
// @Success 200 {array} models.ActiveTab
// @Router /api/in-memory [get]
func (inMemory *InMemory) getActiveTabs(writer http.ResponseWriter, request *http.Request) {

	result := inMemory.activeTabsRepository.GetAll()
	models.ActiveTabs(result).ToJSON(writer)
}

// @Summary Gets Active Tab by Key
// @Produce json
// @Success 200 models.ActiveTab
// @Failure 404 not found
// @Failure 500 internel server error
// @Router /api/in-memory?key={key} [get]
func (inMemory *InMemory) getActiveTab(key string, writer http.ResponseWriter, request *http.Request) {

	activeTabResult, err := inMemory.activeTabsRepository.Get(key)
	if err != nil {
		if err == data.ErrActiveTabNotFound {
			inMemory.l.Printf("Active tab not found for %s %v", key, err.Error())
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		} else {
			inMemory.l.Printf("While getting active tab someting went wrong, Key: %s , Err: %v", key, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := activeTabResult.ToJSON(writer); err != nil {
		inMemory.l.Printf("While converting to json someting went wrong, Key: %s, Err: %v", key, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Adds Active Tab
// @Produce json
// @Success 200 models.ActiveTab
// @Failure 400 Bad Request
// @Failure 500 Internel Server error
// @Router /api/in-memory [post]
func (inMemory *InMemory) addActiveTab(writer http.ResponseWriter, request *http.Request) {

	newActiveTab := &models.ActiveTab{}

	if err := newActiveTab.FromJSON(request.Body); err != nil {
		inMemory.l.Printf("While converting json to struct something went wrong, Err: %v", err)
		http.Error(writer, "Payload format not valid", http.StatusBadRequest)
		return
	}

	if err := inMemory.activeTabsRepository.Add(*newActiveTab); err != nil {
		inMemory.l.Printf("While adding new active tab something went wrong, Err: %v", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	if err := newActiveTab.ToJSON(writer); err != nil {
		inMemory.l.Printf("While converting to json someting went wrong, Err: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
