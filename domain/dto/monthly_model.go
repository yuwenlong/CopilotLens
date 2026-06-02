package dto

type ModelUsage struct {
	Model string  `json:"model"`
	Total float64 `json:"total"`
	Cost  float64 `json:"cost"`
}

type MonthlyModelResponse struct {
	Month  string       `json:"month"`
	Models []ModelUsage `json:"models"`
}
