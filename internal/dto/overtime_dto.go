package dto

type SubmitOvertimeRequest struct {
	Date  string  `json:"date" binding:"required"` // e.g. "2025-06-20"
	Hours float32 `json:"hours" binding:"required,gt=0,lte=3"`
}
