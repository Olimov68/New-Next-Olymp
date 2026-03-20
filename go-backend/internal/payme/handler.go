package payme

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/config"
	"github.com/nextolympservice/go-backend/internal/models"
	"gorm.io/gorm"
)

// ─── JSON-RPC types ──────────────────────────────────────────────────────────

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int               `json:"code"`
	Message map[string]string `json:"message"`
	Data    string            `json:"data,omitempty"`
}

// ─── Param types ─────────────────────────────────────────────────────────────

type CheckPerformParams struct {
	Amount  int64                  `json:"amount"`
	Account map[string]interface{} `json:"account"`
}

type CreateTransactionParams struct {
	ID      string                 `json:"id"`
	Time    int64                  `json:"time"`
	Amount  int64                  `json:"amount"`
	Account map[string]interface{} `json:"account"`
}

type PerformTransactionParams struct {
	ID string `json:"id"`
}

type CancelTransactionParams struct {
	ID     string `json:"id"`
	Reason int    `json:"reason"`
}

type CheckTransactionParams struct {
	ID string `json:"id"`
}

type GetStatementParams struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

// ─── Handler ─────────────────────────────────────────────────────────────────

type Handler struct {
	db  *gorm.DB
	cfg *config.PaymeConfig
}

func NewHandler(db *gorm.DB, cfg *config.PaymeConfig) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// Handle — main JSON-RPC endpoint that Payme calls
func (h *Handler) Handle(c *gin.Context) {
	// 1. Verify Basic Auth
	if !h.verifyAuth(c) {
		return
	}

	// 2. Parse JSON-RPC request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.sendError(c, nil, models.PaymeRPCParseError, "Could not parse request body", "")
		return
	}

	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		h.sendError(c, nil, models.PaymeRPCParseError, "Invalid JSON-RPC request", "")
		return
	}

	log.Printf("[Payme] Method: %s, ID: %v", req.Method, req.ID)

	// 3. Route to method handler
	switch req.Method {
	case "CheckPerformTransaction":
		h.checkPerformTransaction(c, &req)
	case "CreateTransaction":
		h.createTransaction(c, &req)
	case "PerformTransaction":
		h.performTransaction(c, &req)
	case "CancelTransaction":
		h.cancelTransaction(c, &req)
	case "CheckTransaction":
		h.checkTransaction(c, &req)
	case "GetStatement":
		h.getStatement(c, &req)
	default:
		h.sendError(c, req.ID, models.PaymeMethodNotFound, "Method not found", req.Method)
	}
}

// ─── Auth verification ───────────────────────────────────────────────────────

func (h *Handler) verifyAuth(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		h.sendError(c, nil, models.PaymeRPCInternalError, "Authentication required", "auth")
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	if err != nil {
		h.sendError(c, nil, models.PaymeRPCInternalError, "Invalid auth header", "auth")
		return false
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 || parts[0] != "Paycom" {
		h.sendError(c, nil, models.PaymeRPCInternalError, "Invalid credentials", "auth")
		return false
	}

	// Check password against the appropriate key
	expectedKey := h.cfg.Key
	if h.cfg.TestMode {
		expectedKey = h.cfg.TestKey
	}

	if parts[1] != expectedKey {
		h.sendError(c, nil, models.PaymeRPCInternalError, "Invalid credentials", "auth")
		return false
	}

	return true
}

// ─── CheckPerformTransaction ─────────────────────────────────────────────────

func (h *Handler) checkPerformTransaction(c *gin.Context, req *JSONRPCRequest) {
	var params CheckPerformParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendError(c, req.ID, models.PaymeRPCParseError, "Invalid params", "params")
		return
	}

	orderID, ok := h.getOrderID(params.Account)
	if !ok {
		h.sendError(c, req.ID, models.PaymeOrderNotFound,
			"Order not found", "order_id")
		return
	}

	var payment models.Payment
	if err := h.db.First(&payment, orderID).Error; err != nil {
		h.sendError(c, req.ID, models.PaymeOrderNotFound,
			"Order not found", "order_id")
		return
	}

	// Check payment is still pending
	if payment.Status != models.PaymentStatusPending {
		h.sendError(c, req.ID, models.PaymeOrderAlreadyPaid,
			"Order already paid or cancelled", "status")
		return
	}

	// Check amount matches (payment.Amount is UZS, params.Amount is tiyin)
	expectedTiyin := int64(payment.Amount * 100)
	if params.Amount != expectedTiyin {
		h.sendError(c, req.ID, models.PaymeAmountMismatch,
			"Incorrect amount", "amount")
		return
	}

	h.sendResult(c, req.ID, gin.H{
		"allow": true,
	})
}

// ─── CreateTransaction ───────────────────────────────────────────────────────

func (h *Handler) createTransaction(c *gin.Context, req *JSONRPCRequest) {
	var params CreateTransactionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendError(c, req.ID, models.PaymeRPCParseError, "Invalid params", "params")
		return
	}

	orderID, ok := h.getOrderID(params.Account)
	if !ok {
		h.sendError(c, req.ID, models.PaymeOrderNotFound,
			"Order not found", "order_id")
		return
	}

	// Check if this Payme transaction already exists
	var existing models.PaymeTransaction
	if err := h.db.Where("payme_id = ?", params.ID).First(&existing).Error; err == nil {
		// Transaction already exists - check if it's still valid
		if existing.State == models.PaymeStateCreated {
			// Check timeout
			if h.isTimedOut(existing.CreateTime) {
				// Cancel due to timeout
				h.db.Model(&existing).Updates(map[string]interface{}{
					"state":       models.PaymeStateCancelledBefore,
					"cancel_time": time.Now().UnixMilli(),
					"reason":      4, // timeout
				})
				h.sendError(c, req.ID, models.PaymeTransactionTimeout,
					"Transaction timed out", "timeout")
				return
			}
			// Return existing transaction
			h.sendResult(c, req.ID, gin.H{
				"create_time": existing.CreateTime,
				"transaction": fmt.Sprintf("%d", existing.ID),
				"state":       existing.State,
			})
			return
		}
		// Transaction in terminal state
		h.sendError(c, req.ID, models.PaymeCantPerform,
			"Transaction cannot be performed", "state")
		return
	}

	// Find payment
	var payment models.Payment
	if err := h.db.First(&payment, orderID).Error; err != nil {
		h.sendError(c, req.ID, models.PaymeOrderNotFound,
			"Order not found", "order_id")
		return
	}

	if payment.Status != models.PaymentStatusPending {
		h.sendError(c, req.ID, models.PaymeOrderAlreadyPaid,
			"Order already paid or cancelled", "status")
		return
	}

	// Check amount
	expectedTiyin := int64(payment.Amount * 100)
	if params.Amount != expectedTiyin {
		h.sendError(c, req.ID, models.PaymeAmountMismatch,
			"Incorrect amount", "amount")
		return
	}

	// Check if another Payme transaction already exists for this payment
	var otherTx models.PaymeTransaction
	if err := h.db.Where("payment_id = ? AND state = ?", payment.ID, models.PaymeStateCreated).First(&otherTx).Error; err == nil {
		// Another transaction exists for this payment - cancel it
		h.db.Model(&otherTx).Updates(map[string]interface{}{
			"state":       models.PaymeStateCancelledBefore,
			"cancel_time": time.Now().UnixMilli(),
			"reason":      4,
		})
	}

	// Create new PaymeTransaction
	tx := models.PaymeTransaction{
		PaymeID:    params.ID,
		PaymentID:  payment.ID,
		UserID:     payment.UserID,
		Amount:     params.Amount,
		State:      models.PaymeStateCreated,
		CreateTime: params.Time,
	}

	if err := h.db.Create(&tx).Error; err != nil {
		log.Printf("[Payme] Failed to create transaction: %v", err)
		h.sendError(c, req.ID, models.PaymeRPCInternalError,
			"Internal error", "db")
		return
	}

	h.sendResult(c, req.ID, gin.H{
		"create_time": tx.CreateTime,
		"transaction": fmt.Sprintf("%d", tx.ID),
		"state":       tx.State,
	})
}

// ─── PerformTransaction ──────────────────────────────────────────────────────

func (h *Handler) performTransaction(c *gin.Context, req *JSONRPCRequest) {
	var params PerformTransactionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendError(c, req.ID, models.PaymeRPCParseError, "Invalid params", "params")
		return
	}

	var tx models.PaymeTransaction
	if err := h.db.Where("payme_id = ?", params.ID).First(&tx).Error; err != nil {
		h.sendError(c, req.ID, models.PaymeTransactionNotFound,
			"Transaction not found", "id")
		return
	}

	if tx.State == models.PaymeStatePerformed {
		// Already performed - return success (idempotent)
		h.sendResult(c, req.ID, gin.H{
			"transaction":  fmt.Sprintf("%d", tx.ID),
			"perform_time": tx.PerformTime,
			"state":        tx.State,
		})
		return
	}

	if tx.State != models.PaymeStateCreated {
		h.sendError(c, req.ID, models.PaymeCantPerform,
			"Transaction cannot be performed", "state")
		return
	}

	// Check timeout
	if h.isTimedOut(tx.CreateTime) {
		h.db.Model(&tx).Updates(map[string]interface{}{
			"state":       models.PaymeStateCancelledBefore,
			"cancel_time": time.Now().UnixMilli(),
			"reason":      4,
		})
		h.sendError(c, req.ID, models.PaymeCantPerform,
			"Transaction timed out", "timeout")
		return
	}

	performTime := time.Now().UnixMilli()

	// Use DB transaction for atomicity
	err := h.db.Transaction(func(dbTx *gorm.DB) error {
		// Update PaymeTransaction state
		if err := dbTx.Model(&tx).Updates(map[string]interface{}{
			"state":        models.PaymeStatePerformed,
			"perform_time": performTime,
		}).Error; err != nil {
			return err
		}

		// Mark Payment as completed
		if err := dbTx.Model(&models.Payment{}).Where("id = ?", tx.PaymentID).Updates(map[string]interface{}{
			"status":         models.PaymentStatusCompleted,
			"transaction_id": fmt.Sprintf("PAYME-%s", tx.PaymeID),
		}).Error; err != nil {
			return err
		}

		// Credit user balance
		var balance models.Balance
		if err := dbTx.Where("user_id = ?", tx.UserID).First(&balance).Error; err != nil {
			// Create balance if not exists
			balance = models.Balance{
				UserID:   tx.UserID,
				Amount:   0,
				Currency: "UZS",
			}
			dbTx.Create(&balance)
		}

		amountUZS := float64(tx.Amount) / 100.0
		balanceBefore := balance.Amount
		balanceAfter := balanceBefore + amountUZS

		if err := dbTx.Model(&balance).Update("amount", balanceAfter).Error; err != nil {
			return err
		}

		// Create balance transaction record
		balanceTx := models.BalanceTransaction{
			UserID:          tx.UserID,
			Type:            "topup",
			Amount:          amountUZS,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    balanceAfter,
			Description:     "Payme orqali to'ldirish",
			SourceType:      "payment",
			SourceID:        &tx.PaymentID,
			PerformedByType: "system",
		}
		return dbTx.Create(&balanceTx).Error
	})

	if err != nil {
		log.Printf("[Payme] PerformTransaction DB error: %v", err)
		h.sendError(c, req.ID, models.PaymeRPCInternalError,
			"Internal error", "db")
		return
	}

	h.sendResult(c, req.ID, gin.H{
		"transaction":  fmt.Sprintf("%d", tx.ID),
		"perform_time": performTime,
		"state":        models.PaymeStatePerformed,
	})
}

// ─── CancelTransaction ──────────────────────────────────────────────────────

func (h *Handler) cancelTransaction(c *gin.Context, req *JSONRPCRequest) {
	var params CancelTransactionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendError(c, req.ID, models.PaymeRPCParseError, "Invalid params", "params")
		return
	}

	var tx models.PaymeTransaction
	if err := h.db.Where("payme_id = ?", params.ID).First(&tx).Error; err != nil {
		h.sendError(c, req.ID, models.PaymeTransactionNotFound,
			"Transaction not found", "id")
		return
	}

	// Already cancelled - return success (idempotent)
	if tx.State == models.PaymeStateCancelledBefore || tx.State == models.PaymeStateCancelledAfter {
		h.sendResult(c, req.ID, gin.H{
			"transaction": fmt.Sprintf("%d", tx.ID),
			"cancel_time": tx.CancelTime,
			"state":       tx.State,
		})
		return
	}

	cancelTime := time.Now().UnixMilli()
	reason := params.Reason

	if tx.State == models.PaymeStateCreated {
		// Cancel before perform
		err := h.db.Transaction(func(dbTx *gorm.DB) error {
			if err := dbTx.Model(&tx).Updates(map[string]interface{}{
				"state":       models.PaymeStateCancelledBefore,
				"cancel_time": cancelTime,
				"reason":      reason,
			}).Error; err != nil {
				return err
			}
			// Mark payment as failed
			return dbTx.Model(&models.Payment{}).Where("id = ?", tx.PaymentID).Update("status", models.PaymentStatusFailed).Error
		})
		if err != nil {
			log.Printf("[Payme] CancelTransaction DB error: %v", err)
			h.sendError(c, req.ID, models.PaymeRPCInternalError, "Internal error", "db")
			return
		}

		h.sendResult(c, req.ID, gin.H{
			"transaction": fmt.Sprintf("%d", tx.ID),
			"cancel_time": cancelTime,
			"state":       models.PaymeStateCancelledBefore,
		})
		return
	}

	if tx.State == models.PaymeStatePerformed {
		// Cancel after perform - need to refund
		err := h.db.Transaction(func(dbTx *gorm.DB) error {
			if err := dbTx.Model(&tx).Updates(map[string]interface{}{
				"state":       models.PaymeStateCancelledAfter,
				"cancel_time": cancelTime,
				"reason":      reason,
			}).Error; err != nil {
				return err
			}

			// Mark payment as refunded
			if err := dbTx.Model(&models.Payment{}).Where("id = ?", tx.PaymentID).Update("status", models.PaymentStatusRefunded).Error; err != nil {
				return err
			}

			// Deduct from user balance
			var balance models.Balance
			if err := dbTx.Where("user_id = ?", tx.UserID).First(&balance).Error; err != nil {
				return err
			}

			amountUZS := float64(tx.Amount) / 100.0
			balanceBefore := balance.Amount
			balanceAfter := balanceBefore - amountUZS
			if balanceAfter < 0 {
				balanceAfter = 0
			}

			if err := dbTx.Model(&balance).Update("amount", balanceAfter).Error; err != nil {
				return err
			}

			// Create refund balance transaction
			balanceTx := models.BalanceTransaction{
				UserID:          tx.UserID,
				Type:            "refund",
				Amount:          -amountUZS,
				BalanceBefore:   balanceBefore,
				BalanceAfter:    balanceAfter,
				Description:     "Payme to'lov bekor qilindi",
				SourceType:      "refund",
				SourceID:        &tx.PaymentID,
				PerformedByType: "system",
			}
			return dbTx.Create(&balanceTx).Error
		})

		if err != nil {
			log.Printf("[Payme] CancelTransaction (refund) DB error: %v", err)
			h.sendError(c, req.ID, models.PaymeRPCInternalError, "Internal error", "db")
			return
		}

		h.sendResult(c, req.ID, gin.H{
			"transaction": fmt.Sprintf("%d", tx.ID),
			"cancel_time": cancelTime,
			"state":       models.PaymeStateCancelledAfter,
		})
		return
	}

	h.sendError(c, req.ID, models.PaymeCantCancel,
		"Cannot cancel transaction", "state")
}

// ─── CheckTransaction ────────────────────────────────────────────────────────

func (h *Handler) checkTransaction(c *gin.Context, req *JSONRPCRequest) {
	var params CheckTransactionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendError(c, req.ID, models.PaymeRPCParseError, "Invalid params", "params")
		return
	}

	var tx models.PaymeTransaction
	if err := h.db.Where("payme_id = ?", params.ID).First(&tx).Error; err != nil {
		h.sendError(c, req.ID, models.PaymeTransactionNotFound,
			"Transaction not found", "id")
		return
	}

	result := gin.H{
		"create_time":  tx.CreateTime,
		"perform_time": tx.PerformTime,
		"cancel_time":  tx.CancelTime,
		"transaction":  fmt.Sprintf("%d", tx.ID),
		"state":        tx.State,
	}
	if tx.Reason != nil {
		result["reason"] = *tx.Reason
	} else {
		result["reason"] = nil
	}

	h.sendResult(c, req.ID, result)
}

// ─── GetStatement ────────────────────────────────────────────────────────────

func (h *Handler) getStatement(c *gin.Context, req *JSONRPCRequest) {
	var params GetStatementParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendError(c, req.ID, models.PaymeRPCParseError, "Invalid params", "params")
		return
	}

	var transactions []models.PaymeTransaction
	h.db.Where("create_time BETWEEN ? AND ?", params.From, params.To).
		Order("create_time ASC").
		Find(&transactions)

	var result []gin.H
	for _, tx := range transactions {
		item := gin.H{
			"id":           tx.PaymeID,
			"time":         tx.CreateTime,
			"amount":       tx.Amount,
			"account":      gin.H{"order_id": fmt.Sprintf("%d", tx.PaymentID)},
			"create_time":  tx.CreateTime,
			"perform_time": tx.PerformTime,
			"cancel_time":  tx.CancelTime,
			"transaction":  fmt.Sprintf("%d", tx.ID),
			"state":        tx.State,
		}
		if tx.Reason != nil {
			item["reason"] = *tx.Reason
		} else {
			item["reason"] = nil
		}
		result = append(result, item)
	}

	if result == nil {
		result = []gin.H{}
	}

	h.sendResult(c, req.ID, gin.H{
		"transactions": result,
	})
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func (h *Handler) getOrderID(account map[string]interface{}) (uint, bool) {
	orderIDRaw, exists := account["order_id"]
	if !exists {
		return 0, false
	}

	switch v := orderIDRaw.(type) {
	case float64:
		return uint(v), true
	case string:
		var id uint
		if _, err := fmt.Sscanf(v, "%d", &id); err == nil {
			return id, true
		}
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return uint(n), true
		}
	}
	return 0, false
}

func (h *Handler) isTimedOut(createTime int64) bool {
	return time.Now().UnixMilli()-createTime > models.PaymeTimeoutMs
}

func (h *Handler) sendResult(c *gin.Context, id interface{}, result interface{}) {
	c.JSON(http.StatusOK, JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

func (h *Handler) sendError(c *gin.Context, id interface{}, code int, message string, data string) {
	c.JSON(http.StatusOK, JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code: code,
			Message: map[string]string{
				"uz": message,
				"ru": message,
				"en": message,
			},
			Data: data,
		},
	})
}
