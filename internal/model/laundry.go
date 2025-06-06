package model

import "time"

type (
	LaundryResponse struct {
		Id                string    `json:"id" db:"id"`
		DetailNumber      string    `json:"detail_number" db:"detail_number"`
		Title             string    `json:"title" db:"title"`
		LaundryDate       time.Time `json:"-" db:"laundry_date"`
		LaundryDateString string    `json:"laundry_date"`
		TotalItems        int       `json:"total_items" db:"total_items"`
		Status            int       `json:"-" db:"status"`
		StatusLabel       string    `json:"status_label"`
	}

	LaundryQueryParam struct {
		CategoryName    string
		LaundryDateFrom *time.Time
		LaundryDateTo   *time.Time
		DetailNumber    string
		Page            int
	}
)
