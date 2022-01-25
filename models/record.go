package models

import (
	"encoding/json"
	"io"
	"time"
)

// record database entity model
type Record struct {
	Key        string    `json:"key" bson:"key"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	TotalCount int       `json:"totalCount" bson:"totalCount"`
}

// converts JSON to Record struct
func (rd *Record) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(rd)
}

//converts Record to JSON
func (r *Record) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}
