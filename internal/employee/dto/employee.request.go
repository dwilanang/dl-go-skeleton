package dto

type AttendanceRequest struct {
	UserID int64  `json:"user_id" swaggerignore:"true"`
	Date   string `json:"date" binding:"required"`
	By     int64  `json:"by" swaggerignore:"true"`
}

type OvertimeRequest struct {
	Hours  int    `json:"hours" binding:"required"`
	UserID int64  `json:"user_id" swaggerignore:"true"`
	Date   string `json:"date" binding:"required"`
	By     int64  `json:"by" swaggerignore:"true"`
}

type ReimbursementRequest struct {
	Amount      int64  `json:"amount" binding:"required"`
	Description string `json:"description" binding:"required"`
	UserID      int64  `json:"user_id" swaggerignore:"true"`
	Date        string `json:"date" binding:"required"`
	By          int64  `json:"by" swaggerignore:"true"`
}
