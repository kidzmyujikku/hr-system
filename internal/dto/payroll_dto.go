package dto

type PayrollRequests struct {
	PayCycleId uint `json:"paycycle_id" binding:"required"`
}
