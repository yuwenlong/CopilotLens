package dto

type DailyTotal struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type MonthlyTotalResponse struct {
	Month string       `json:"month"`
	Total float64      `json:"total"`
	Daily []DailyTotal `json:"daily"`
}
