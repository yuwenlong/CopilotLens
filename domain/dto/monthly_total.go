package dto

type MonthlyTotalResponse struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}
