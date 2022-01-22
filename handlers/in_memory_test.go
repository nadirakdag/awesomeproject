package handlers

import (
	"awesomeProject/data"
	"awesomeProject/models"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInMemoryServeHTTP(t *testing.T) {

	const path string = "/api/in-memory"

	l := log.Default()
	l.SetOutput(ioutil.Discard)

	inMemory := &InMemory{
		l: l,
		keyValuePairRepository: &data.KeyValueInMemoryRepository{
			KeyValues: []models.KeyValuePair{
				{Key: "test", Value: "nadir"},
			},
		},
	}

	t.Run("returns method not allowed when send other then POST", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(inMemory.ServeHTTP)
		req, err := http.NewRequest("PUT", path, nil)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("return bad request error when it can not parse body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(inMemory.ServeHTTP)

		body := strings.NewReader("a")
		req, err := http.NewRequest("POST", path, body)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("return bad request error when sended key already exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(inMemory.ServeHTTP)

		keyValuePair := &models.KeyValuePair{Key: "test", Value: "nadir"}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(keyValuePair)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", path, &buf)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("return not found error when sended key does not exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(inMemory.ServeHTTP)

		req, err := http.NewRequest("GET", path+"?key=test1", nil)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("return ok and key value pair when sended key successfully added", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(inMemory.ServeHTTP)

		keyValuePair := &models.KeyValuePair{Key: "test1", Value: "nadir"}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(keyValuePair)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", path, &buf)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		returnedKeyValuePair := &models.KeyValuePair{}
		if err := json.NewDecoder(rr.Body).Decode(returnedKeyValuePair); err != nil {
			t.Fatalf("body decode went wrong, body %v", rr.Body.String())
		}

		if keyValuePair.Key != returnedKeyValuePair.Key || keyValuePair.Value != returnedKeyValuePair.Value {
			t.Errorf("handler returned wrong key value pair, got %s want %v", rr.Body.String(), keyValuePair)
		}
	})

	t.Run("return ok and all key pairs", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(inMemory.ServeHTTP)

		req, err := http.NewRequest("GET", path, nil)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
