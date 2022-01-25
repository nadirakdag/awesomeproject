package models

import (
	"encoding/json"
	"io"
)

// this model used for filtering records
// also used for records POST request model
type RecordFilter struct {
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
	MinCount  int    `json:"minCount" validate:"required,numeric"`
	MaxCount  int    `json:"maxCount" validate:"required,numeric"`
}

// converts JSON to RecordFilter struct
func (rf *RecordFilter) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(rf)
}

//converts Record to JSON
func (r *RecordFilter) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}
