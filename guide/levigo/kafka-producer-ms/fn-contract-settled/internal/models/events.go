package models

import "time"

type ContractSettled struct {
	ContractID  string
	Document    string
	TotalAmount float64
	SettledAt   time.Time
}
