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
	d := json.NewDecoder(r)
	return d.Decode(activeTab)
}

func (activeTab *ActiveTab) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(activeTab)
}

type ActiveTabs []ActiveTab

func (activeTabs ActiveTabs) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(activeTabs)
}
