package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Handler struct {
	orderHandler *OrderHandler
	logger       *zap.Logger
}

func NewHandler(orderHandler *OrderHandler, logger *zap.Logger) *Handler {
	return &Handler{
		orderHandler: orderHandler,
		logger:       logger,
	}
}
func (h *Handler) InitRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Route("/order", func(r chi.Router) {
		r.Route("/orders", func(r chi.Router) {
			r.Get("/{id}/status", h.orderHandler.GetStatus)

			r.Group(func(r chi.Router) {
				r.Use(AuthMiddleware)
				r.Post("/", h.orderHandler.CreateOrder)
			})
		})
	})
	return r
}
