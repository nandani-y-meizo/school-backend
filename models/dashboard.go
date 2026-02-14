package models

import "time"

type DashboardStats struct {
	Collection CollectionStats `json:"collection"`
	FeesStatus FeesStatusStats `json:"fees_status"`
	Holidays   []Holiday       `json:"holidays"`
}

type CollectionStats struct {
	TotalPaidAmount float64 `json:"total_paid_amount"`
	CashAmount      float64 `json:"cash_amount"`
	UPIAmount       float64 `json:"upi_amount"`
}

type FeesStatusStats struct {
	PaidStudents   int `json:"paid_students"`
	UnpaidStudents int `json:"unpaid_students"`
	TotalStudents  int `json:"total_students"`
}

type Holiday struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
	Flag string    `json:"flag"`
}
