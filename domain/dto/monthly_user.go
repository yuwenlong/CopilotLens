package dto

type UserModel struct {
	Model string  `json:"model"`
	Total float64 `json:"total"`
	Cost  float64 `json:"cost"`
}

type UserUsage struct {
	Username string      `json:"username"`
	Name     string      `json:"name"`
	Total    float64     `json:"total"`
	Cost     float64     `json:"cost"`
	Models   []UserModel `json:"models"`
}

type MonthlyUserResponse struct {
	Month string      `json:"month"`
	Users []UserUsage `json:"users"`
}
