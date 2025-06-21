package dto

type PayCycleRequests struct {
	StartDate string `json:"start_date" binding:"required"` // "2025-06-01"
	EndDate   string `json:"end_date" binding:"required"`   // "2025-06-30"
}
