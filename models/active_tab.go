package models

import (
	"encoding/json"
	"io"
)

type ActiveTab struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (activeTab *ActiveTab) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(activeTab)
}

func (activeTab *ActiveTab) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(activeTab)
}

type ActiveTabs []ActiveTab

func (activeTabs ActiveTabs) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(activeTabs)
}
