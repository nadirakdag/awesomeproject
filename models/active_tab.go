package models

import (
	"encoding/json"
	"io"
)

// activetab entity model
type ActiveTab struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// converts JSON to activeTab struct
func (activeTab *ActiveTab) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(activeTab)
}

//converts activeTab to JSON
func (activeTab *ActiveTab) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(activeTab)
}

// activeTab Array model
type ActiveTabs []ActiveTab

// converts activeTabs array to JSON
func (activeTabs ActiveTabs) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(activeTabs)
}
