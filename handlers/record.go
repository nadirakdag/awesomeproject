package handlers

import (
	"awesomeProject/data"
	"awesomeProject/helpers"
	"awesomeProject/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

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

var ErrStartDateFormatInvalid = errors.New("start date format is invalid, is should be YYYY-MM-DD")
var ErrEndDateFormatInvalid = errors.New("end date format is invalid, is should be YYYY-MM-DD")

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

	if err := json.NewDecoder(request.Body).Decode(filter); err != nil {
		r.l.Printf("Error while decode filter json %v \n", err)
		result := errorResult(-1, err.Error())
		returnJSONError(writer, result, 500)
		return
	}

	if err := validator.New().Struct(filter); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		r.l.Printf("json validation error %v \n", err)
		result := errorResult(-1, validationErrors.Error())
		returnJSONError(writer, result, 400)
		return
	}

	if _, err := time.Parse(helpers.DateFormat, filter.StartDate); err != nil {
		result := errorResult(-1, err.Error())
		returnJSONError(writer, result, 400)
		return
	}

	if _, err := time.Parse(helpers.DateFormat, filter.EndDate); err != nil {
		result := errorResult(-1, err.Error())
		returnJSONError(writer, result, 400)
		return
	}

	records, err := r.repository.Get(filter)
	if err != nil {
		r.l.Printf("Error whilte getting records, %v \n", err)
		result := errorResult(-2, err.Error())
		returnJSONError(writer, result, 500)
		return
	}

	result := &models.RecordResult{
		Code:    0,
		Message: "Success",
		Records: records,
	}

	returnJSON(writer, result)
}

// creates Record result for errors
func errorResult(code int, msg string) *models.RecordResult {
	return &models.RecordResult{
		Code:    code,
		Message: msg,
		Records: nil,
	}
}
