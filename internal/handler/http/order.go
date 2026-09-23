package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AndroDeMohawk/flash-sales/internal/infrastructure/redis"
	"github.com/AndroDeMohawk/flash-sales/internal/usecase"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type OrderHandler struct {
	orderUC *usecase.OrderUseCase
	logger  *zap.Logger
}

func NewOrderHandler(orderUC *usecase.OrderUseCase, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderUC: orderUC,
		logger:  logger,
	}
}

type CreateOrderRequest struct {
	UserID   int64 `json:"user_id"`
	TicketID int64 `json:"ticket_id"`
	Quantity int   `json:"quantity"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// CreateOrder — POST /api/v1/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("UserIdCtxKey").(int64)
	if !ok {
		http.Error(w, "Missing UserID", http.StatusBadRequest)
		return
	}
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 || req.TicketID <= 0 || req.Quantity <= 0 {
		http.Error(w, "user_id, ticket_id and quantity must be positive", http.StatusBadRequest)
		return
	}

	orderID, err := h.orderUC.CreateOrder(r.Context(), userId, req.TicketID, req.Quantity)
	if err != nil {
		if errors.Is(err, redis.ErrSoldOut) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // 409 Sold Out
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "билеты закончились",
			})
			return
		}

		h.logger.Error("failed to process order creation", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 Accepted
	_ = json.NewEncoder(w).Encode(CreateOrderResponse{
		OrderID: orderID,
		Status:  "PENDING",
		Message: "Заказ зарезервирован и отправлен в обработку",
	})
}

// GetStatus — GET /api/v1/orders/{id}/status
func (h *OrderHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		http.Error(w, "missing order id parameter", http.StatusBadRequest)
		return
	}

	status, err := h.orderUC.GetOrderStatus(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order status not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"order_id": orderID,
		"status":   status,
	})
}
