package handlers

import (
	"awesomeProject/data"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator"
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
	if err := j.Decode(filter); err != nil {
		r.l.Printf("Error while decode filter json %v \n", err)
		result := errorResult(-1, err.Error())
		jsonError(writer, result, 500)
		return
	}

	if err := validator.New().Struct(filter); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		r.l.Printf("json validation error %v \n", err)
		http.Error(writer, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	records, err := r.repository.Get(filter)
	if err != nil {
		r.l.Printf("Error whilte getting records, %v \n", err)

		if err == data.ErrEndDateFormatInvalid || err == data.ErrStartDateFormatInvalid {
			result := errorResult(-1, err.Error())
			jsonError(writer, result, 400)
			return
		} else {
			result := errorResult(-2, err.Error())
			jsonError(writer, result, 500)
			return
		}
	}

	result := &RecordResult{
		Code:    0,
		Message: "Success",
		Records: records,
	}

	json.NewEncoder(writer).Encode(result)
}

func jsonError(writer http.ResponseWriter, content interface{}, code int) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(code)
	json.NewEncoder(writer).Encode(content)
}

func errorResult(code int, msg string) *RecordResult {
	return &RecordResult{
		Code:    code,
		Message: msg,
		Records: nil,
	}
}
