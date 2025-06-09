package model

import "time"

type (
	LaundryResponse struct {
		Id                string    `json:"id" db:"id"`
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
		Page            int
	}
)

type (
	CategoryRequest struct {
		Name string `json:"name" validate:"required"`
	}

	CategoryResponse struct {
		Id       string `json:"id" db:"id"`
		Name     string `json:"name" db:"name"`
		UserId   string `json:"-" db:"user_id"`
		IsActive bool   `json:"-" db:"is_active"`
	}
)

type (
	AddLaundryRequest struct {
		Title       string                `json:"title" validate:"required"`
		LaundryDate string                `json:"laundry_date" validate:"required,date_format"`
		Items       []LaundryItemsRequest `json:"items" validate:"required,gt=1"`
	}

	LaundryItemsRequest struct {
		CategoryId string  `json:"category_id" validate:"required"`
		Amount     int     `json:"amount" validate:"required,min=1"`
		Notes      *string `json:"notes"`
	}
)
