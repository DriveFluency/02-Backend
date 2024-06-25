package domain

import (
	//"time"
)

type Pay struct {
	Id      int       `json:"id"`
	Date    string `json:"date" ` //,omitempty
	Method  string    `json:"method" binding:"required"`
	Amount  float64   `json:"amount" binding:"required"`
	Receipt string    `json:"receipt" binding:"required"`
}
