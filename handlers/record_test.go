package handlers

import (
	"awesomeProject/models"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServeHTTP(t *testing.T) {

	const method string = "POST"
	const path string = "/api/records"

	l := log.Default()
	l.SetOutput(ioutil.Discard)

	record := &Records{
		l: l,
	}

	t.Run("returns method not allowed when send other then POST", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(record.ServeHTTP)
		req, err := http.NewRequest("PUT", path, nil)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("return internal server error when it can not parse body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(record.ServeHTTP)
		body := strings.NewReader("a")
		req, err := http.NewRequest(method, path, body)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		checkResponseCode(rr, t, -1)
	})

	t.Run("returns bad request error when validation failed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(record.ServeHTTP)
		recordFilter := &models.RecordFilter{}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(recordFilter)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest(method, path, &buf)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		checkResponseCode(rr, t, -1)
	})

	t.Run("returns bad request when start or end date format wrong", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(record.ServeHTTP)
		recordFilter := &models.RecordFilter{StartDate: "12-01-2021", EndDate: "2022-01-12", MinCount: 1, MaxCount: 1}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(recordFilter)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest(method, path, &buf)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		checkResponseCode(rr, t, -1)
	})

	t.Run("returns success when body valid", func(t *testing.T) {

		record.repository = &MockSucessRecordRepository{}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(record.ServeHTTP)
		recordFilter := &models.RecordFilter{StartDate: "2022-01-12", EndDate: "2022-01-12", MinCount: 1, MaxCount: 1}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(recordFilter)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest(method, path, &buf)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		checkResponseCode(rr, t, 0)
	})

	t.Run("returns internal server error when repository returns an error", func(t *testing.T) {

		record.repository = &MockFailRecordRepository{}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(record.ServeHTTP)
		recordFilter := &models.RecordFilter{StartDate: "2022-01-12", EndDate: "2022-01-12", MinCount: 1, MaxCount: 1}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(recordFilter)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest(method, path, &buf)

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		checkResponseCode(rr, t, -2)
	})
}

func checkResponseCode(rr *httptest.ResponseRecorder, t *testing.T, want int) {
	result := &RecordResult{}
	_ = json.NewDecoder(rr.Body).Decode(result)

	if result.Code != want {
		t.Errorf("handler returned wrong Response Model Code: got %v want %v", result.Code, want)
	}
}

type MockSucessRecordRepository struct {
}

func (mock *MockSucessRecordRepository) Get(filter *models.RecordFilter) ([]models.Record, error) {
	return []models.Record{
		{Key: "TEST", CreatedAt: time.Now(), TotalCount: 300},
	}, nil
}

type MockFailRecordRepository struct {
}

func (mock *MockFailRecordRepository) Get(filter *models.RecordFilter) ([]models.Record, error) {
	return nil, errors.New("for test purpose")
}
