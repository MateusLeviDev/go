package models

import "time"

type Contract struct {
	ID          string
	Document    string
	TotalAmount float64
	SettledAt   time.Time
}
