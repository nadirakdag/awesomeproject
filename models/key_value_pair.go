package models

import (
	"encoding/json"
	"io"
)

// KeyValuePair entity model
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type JSONHelper interface {
	FromJSON(r io.Reader) error
	ToJSON(w io.Writer) error
}

// converts JSON to KeyValuePair struct
func (keyValuePair *KeyValuePair) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(keyValuePair)
}

//converts KeyValuePair to JSON
func (keyValuePair *KeyValuePair) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(keyValuePair)
}

// KeyValuePair Array model
type KeyValuePairs []KeyValuePair

// converts KeyValuePairs array to JSON
func (keyValuePairs KeyValuePairs) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(keyValuePairs)
}

// converts JSON to KeyValuePairs struct
func (keyValuePairs KeyValuePairs) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(keyValuePairs)
}
