package payments

// ListPaymentsParams — to'lovlar filtrlash
type ListPaymentsParams struct {
	UserID     uint   `form:"user_id"`
	SourceType string `form:"source_type" binding:"omitempty,oneof=olympiad mock_test"`
	SourceID   uint   `form:"source_id"`
	Status     string `form:"status" binding:"omitempty,oneof=pending completed failed refunded"`
	Search     string `form:"search"`
	DateFrom   string `form:"date_from"` // YYYY-MM-DD
	DateTo     string `form:"date_to"`
	SortBy     string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at amount status"`
	SortOrder  string `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	Limit      int    `form:"limit,default=20" binding:"min=1,max=100"`
}

// UpdatePaymentStatusRequest — to'lov statusini o'zgartirish
type UpdatePaymentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed failed refunded"`
	Note   string `json:"note" binding:"max=500"`
}

// ManualPaymentRequest — qo'lda to'lov yaratish
type ManualPaymentRequest struct {
	UserID     uint    `json:"user_id" binding:"required,min=1"`
	SourceType string  `json:"source_type" binding:"required,oneof=olympiad mock_test"`
	SourceID   uint    `json:"source_id" binding:"required,min=1"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Note       string  `json:"note" binding:"max=500"`
}

// RefundRequest — qaytarish so'rovi
type RefundRequest struct {
	Reason string `json:"reason" binding:"required,min=5,max=500"`
}
