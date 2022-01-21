package handlers

import (
	"awesomeProject/data"
	"encoding/json"
	"log"
	"net/http"
)

type Records struct {
	l          *log.Logger
	repository data.RecordRepository
}

func NewRecord(l *log.Logger, repository data.RecordRepository) *Records {
	return &Records{
		l:          l,
		repository: repository,
	}
}

type RecordResult struct {
	Code    int           `json:"code"`
	Message string        `json:"msg"`
	Records []data.Record `json:"records"`
}

func (r *Records) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	r.l.Println("Handle Records requested")

	filter := &data.RecordFilter{}
	j := json.NewDecoder(request.Body)
	err := j.Decode(filter)

	if err != nil {
		http.Error(writer, "Payload format not valid", http.StatusBadRequest)
		return
	}

	records := r.repository.Get(filter)

	result := &RecordResult{
		Code:    0,
		Message: "Success",
		Records: records,
	}

	jE := json.NewEncoder(writer)
	jE.Encode(result)
}
