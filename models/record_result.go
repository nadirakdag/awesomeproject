package models

import (
	"encoding/json"
	"io"
)

// @Summary Record GET Response Struct
type RecordResult struct {
	Code    int      `json:"code"`
	Message string   `json:"msg"`
	Records []Record `json:"records"`
}

// converts JSON to RecordFilter struct
func (rr *RecordResult) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(rr)
}

//converts Record to JSON
func (r *RecordResult) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}
