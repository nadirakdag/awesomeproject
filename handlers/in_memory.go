package handlers

import (
	"awesomeProject/data"
	"awesomeProject/models"
	"log"
	"net/http"
)

// is a http handler object
type InMemory struct {
	l                      *log.Logger
	keyValuePairRepository data.KeyValueRepository
}

// Creates a new InMemory handler with given logger and key value repository
func NewInMemory(l *log.Logger, KeyValuePairRepository data.KeyValueRepository) *InMemory {
	return &InMemory{l: l, keyValuePairRepository: KeyValuePairRepository}
}

// @Summary ServeHTTP is the main entry point for the handler and staisfies the http.Handler
func (inMemory *InMemory) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {
		inMemory.l.Println("InMemory GET", request.URL)
		key := request.URL.Query().Get("key")
		if key == "" {
			inMemory.l.Println("InMemory GET, QueryParam Key is null returning all key value pair ")
			inMemory.getKeyValuePairs(writer, request)
			return
		} else {
			inMemory.l.Println("InMemory GET, QueryParam Key is found returning key value pair")
			inMemory.getKeyValuePair(key, writer, request)
			return
		}
	} else if request.Method == http.MethodPost {
		inMemory.addKeyValuePair(writer, request)
		return
	}

	writer.WriteHeader(http.StatusMethodNotAllowed)
}

// @Summary Get All Key Value Pairs
// @Produce json
// @Success 200 {array} models.KeyValuePair
// @Router /api/in-memory [get]
func (inMemory *InMemory) getKeyValuePairs(writer http.ResponseWriter, request *http.Request) {

	result := inMemory.keyValuePairRepository.GetAll()
	returnJSON(writer, result)
}

// @Summary Gets Key Value Pair by Key
// @Produce json
// @Success 200 models.KeyValuePair
// @Failure 404 not found
// @Failure 500 internel server error
// @Router /api/in-memory?key={key} [get]
func (inMemory *InMemory) getKeyValuePair(key string, writer http.ResponseWriter, request *http.Request) {

	keyValuePairResult, err := inMemory.keyValuePairRepository.Get(key)
	if err != nil {
		if err == data.ErrKeyValuePairNotFound {
			inMemory.l.Printf("Key Value Pair not found for %s %v", key, err.Error())
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		} else {
			inMemory.l.Printf("While getting key value pair someting went wrong, Key: %s , Err: %v", key, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := returnJSON(writer, keyValuePairResult); err != nil {
		inMemory.l.Printf("While converting to json someting went wrong, Key: %s, Err: %v", key, err)
		return
	}
}

// @Summary Adds Key Value Pair
// @Produce json
// @Success 200 models.KeyValuePair
// @Failure 400 Bad Request
// @Failure 500 Internel Server error
// @Router /api/in-memory [post]
func (inMemory *InMemory) addKeyValuePair(writer http.ResponseWriter, request *http.Request) {

	newKeyValuePair := &models.KeyValuePair{}

	if err := newKeyValuePair.FromJSON(request.Body); err != nil {
		inMemory.l.Printf("While converting json to struct something went wrong, Err: %v", err)
		http.Error(writer, "Payload format not valid", http.StatusBadRequest)
		return
	}

	if err := inMemory.keyValuePairRepository.Add(*newKeyValuePair); err != nil {
		inMemory.l.Printf("While adding new key value pair something went wrong, Err: %v", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	if err := returnJSON(writer, newKeyValuePair); err != nil {
		inMemory.l.Printf("While converting to json someting went wrong, Err: %v", err)
		return
	}
}
