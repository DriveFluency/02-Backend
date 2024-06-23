package domain

type Pack struct {
	Id              int     `json:"id"`
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description" binding:"required"`
	NumberClasses   int     `json:"number_classes" binding:"required"`
	DurationClasses int     `json:"duration_classes" binding:"required"`
	Cost            float64 `json:"cost" binding:"required"`
}
