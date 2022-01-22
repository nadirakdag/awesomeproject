package handlers

import (
	"awesomeProject/data"
	"awesomeProject/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

// is a http handler object
type Records struct {
	l          *log.Logger
	repository data.RecordRepository
}

// Creates a new Record handler with given logger and record repository
func NewRecord(l *log.Logger, repository data.RecordRepository) *Records {
	return &Records{
		l:          l,
		repository: repository,
	}
}

// @Summary Record GET Response Struct
type RecordResult struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Records []models.Record `json:"records"`
}

// @Summary ServeHTTP is the main entry point for the handler and staisfies the http.Handler
// @Summary Returns filtered records
// @Produce json
// @Success 200 {object} RecordResult
// @Failure 405 Method not allowed other then POST
// @Failure 500 internal server error {object} RecordResult
// @Failure 400 bad request {object} RecordResult
// @Router /api/recors [post]
func (r *Records) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.l.Println("Handle Records requested")

	filter := &models.RecordFilter{}
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

// returns http error with givven http status code and object
func jsonError(writer http.ResponseWriter, content interface{}, code int) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(code)
	json.NewEncoder(writer).Encode(content)
}

// creates Record result for errors
func errorResult(code int, msg string) *RecordResult {
	return &RecordResult{
		Code:    code,
		Message: msg,
		Records: nil,
	}
}
