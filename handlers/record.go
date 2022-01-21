package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Records struct {
	l *log.Logger
}

func NewRecord(l *log.Logger) *Records {
	return &Records{l: l}
}

type RecordFilter struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type RecordDto struct {
	Key        string `json:"key"`
	CreatedAt  string `json:"createdAt"`
	TotalCount int    `json:"totalCount"`
}

type RecordResult struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Records []RecordDto `json:"records"`
}

func (r *Records) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	r.l.Println("Handle Records requested")

	filter := &RecordFilter{}
	j := json.NewDecoder(request.Body)
	err := j.Decode(filter)

	if err != nil {
		http.Error(writer, "Payload format not valid", http.StatusBadRequest)
		return
	}

	result := &RecordResult{
		Code:    1,
		Message: "Success",
		Records: []RecordDto{
			{
				Key:        "TAKwGc6Jr4i8Z487",
				CreatedAt:  "2017-01-28T01:22:14.398Z",
				TotalCount: 2800,
			},
			{
				Key:        "NAeQ8eX7e5TEg7oH",
				CreatedAt:  "2017-01-27T08:19:14.135Z",
				TotalCount: 2900,
			},
		},
	}

	jE := json.NewEncoder(writer)
	jE.Encode(result)
}
