package models

type RecordFilter struct {
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
	MinCount  int    `json:"minCount" validate:"required,numeric"`
	MaxCount  int    `json:"maxCount" validate:"required,numeric"`
}
