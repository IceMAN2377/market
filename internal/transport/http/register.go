package http

import (
	"github.com/IceMAN2377/market/internal/service"

	"log/slog"
	"net/http"
)

func RegisterEndpoints(logger *slog.Logger, router *http.ServeMux, service service.Service) {
	handler := newHandler(service, logger)

	// CRUD операции для подписок
	router.HandleFunc("POST /api/v1/subscriptions", handler.CreateSubscription)
	router.HandleFunc("GET /api/v1/subscriptions", handler.GetSubscriptions)
	router.HandleFunc("GET /api/v1/subscriptions/{id}", handler.GetSubscription)
	router.HandleFunc("PUT /api/v1/subscriptions/{id}", handler.UpdateSubscription)
	router.HandleFunc("DELETE /api/v1/subscriptions/{id}", handler.DeleteSubscription)

	// Расчет стоимости
	router.HandleFunc("POST /api/v1/subscriptions/cost-calculation", handler.CalculateCost)

	RegisterSwaggerEndpoints(router)
}
