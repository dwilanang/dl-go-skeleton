package dto

type UserResponse struct {
	Data UserData `json:"data"`
}

type UserData struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

type UserSalaryResponse struct {
	FullName string  `json:"full_name"`
	Amount   float64 `json:"amount"`
}
