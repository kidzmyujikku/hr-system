package dto

type SubmitReimbursementRequest struct {
	Amount      int64  `json:"amount" binding:"required,gte=1"`        // Amount in Rupiah (must be > 0)
	Description string `json:"description" binding:"required,max=255"` // Reason for the reimbursement
	Date        string `json:"date" binding:"required"`                // e.g. "2025-06-20"
}
