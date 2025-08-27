package http

import (
	"encoding/json"
	"errors"
	errs "github.com/IceMAN2377/market/internal/errors"
	"github.com/IceMAN2377/market/internal/models"
	"github.com/IceMAN2377/market/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

func newHandler(service service.Service, logger *slog.Logger) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

type handler struct {
	service service.Service
	logger  *slog.Logger
}

// CreateSubscription создает новую подписку
func (h *handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		ResponseWithError(h.logger, w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	// Базовая валидация обязательных полей
	if req.ServiceName == "" {
		ResponseWithError(h.logger, w, "service_name is required", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		ResponseWithError(h.logger, w, "user_id is required", http.StatusBadRequest)
		return
	}
	if req.StartDate == "" {
		ResponseWithError(h.logger, w, "start_date is required", http.StatusBadRequest)
		return
	}
	if req.Price <= 0 {
		ResponseWithError(h.logger, w, "price must be positive", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.CreateSubscription(ctx, &req)
	if err != nil {
		h.logger.Error("failed to create subscription", "error", err)

		if errors.Is(err, errs.ErrInvalidUUID) ||
			errors.Is(err, errs.ErrInvalidDateFormat) ||
			errors.Is(err, errs.ErrInvalidDateRange) ||
			errors.Is(err, errs.ErrInvalidPrice) {
			ResponseWithError(h.logger, w, err.Error(), http.StatusBadRequest)
			return
		}

		ResponseWithError(h.logger, w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription created", "subscription_id", subscription.ID, "user_id", subscription.UserID)
	Response(h.logger, w, subscription, http.StatusCreated)
}

// GetSubscription получает подписку по ID
func (h *handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ResponseWithError(h.logger, w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetSubscriptionByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to get subscription", "error", err, "id", id)

		if errors.Is(err, errs.ErrNotFound) {
			ResponseWithError(h.logger, w, "subscription not found", http.StatusNotFound)
			return
		}

		ResponseWithError(h.logger, w, "internal server error", http.StatusInternalServerError)
		return
	}

	Response(h.logger, w, subscription, http.StatusOK)
}

// GetSubscriptions получает список подписок с фильтрацией
func (h *handler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters := &models.SubscriptionFilters{
		Limit:  10, // значение по умолчанию
		Offset: 0,  // значение по умолчанию
	}

	// Парсинг query параметров
	query := r.URL.Query()

	if userID := query.Get("user_id"); userID != "" {
		filters.UserID = &userID
	}

	if serviceName := query.Get("service_name"); serviceName != "" {
		filters.ServiceName = &serviceName
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	response, err := h.service.GetSubscriptions(ctx, filters)
	if err != nil {
		h.logger.Error("failed to get subscriptions", "error", err)

		if errors.Is(err, errs.ErrInvalidUUID) {
			ResponseWithError(h.logger, w, err.Error(), http.StatusBadRequest)
			return
		}

		ResponseWithError(h.logger, w, "internal server error", http.StatusInternalServerError)
		return
	}

	Response(h.logger, w, response, http.StatusOK)
}

// UpdateSubscription обновляет подписку
func (h *handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ResponseWithError(h.logger, w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		ResponseWithError(h.logger, w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	// Проверяем, что хотя бы одно поле для обновления указано
	if req.ServiceName == nil && req.Price == nil && req.EndDate == nil {
		ResponseWithError(h.logger, w, "at least one field must be provided for update", http.StatusBadRequest)
		return
	}

	// Валидация цены, если она указана
	if req.Price != nil && *req.Price <= 0 {
		ResponseWithError(h.logger, w, "price must be positive", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.UpdateSubscription(ctx, id, &req)
	if err != nil {
		h.logger.Error("failed to update subscription", "error", err, "id", id)

		if errors.Is(err, errs.ErrNotFound) {
			ResponseWithError(h.logger, w, "subscription not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, errs.ErrInvalidDateFormat) ||
			errors.Is(err, errs.ErrInvalidDateRange) ||
			errors.Is(err, errs.ErrInvalidPrice) ||
			errors.Is(err, errs.ErrInvalidData) {
			ResponseWithError(h.logger, w, err.Error(), http.StatusBadRequest)
			return
		}

		ResponseWithError(h.logger, w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription updated", "subscription_id", id)
	Response(h.logger, w, subscription, http.StatusOK)
}

// DeleteSubscription удаляет подписку
func (h *handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ResponseWithError(h.logger, w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSubscription(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete subscription", "error", err, "id", id)

		if errors.Is(err, errs.ErrNotFound) {
			ResponseWithError(h.logger, w, "subscription not found", http.StatusNotFound)
			return
		}

		ResponseWithError(h.logger, w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription deleted", "subscription_id", id)
	w.WriteHeader(http.StatusNoContent)
}

// CalculateCost рассчитывает стоимость подписок за период
func (h *handler) CalculateCost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CostCalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		ResponseWithError(h.logger, w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if req.StartDate == "" {
		ResponseWithError(h.logger, w, "start_date is required", http.StatusBadRequest)
		return
	}
	if req.EndDate == "" {
		ResponseWithError(h.logger, w, "end_date is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.CalculateCost(ctx, &req)
	if err != nil {
		h.logger.Error("failed to calculate cost", "error", err)

		if errors.Is(err, errs.ErrInvalidUUID) ||
			errors.Is(err, errs.ErrInvalidDateFormat) ||
			errors.Is(err, errs.ErrInvalidDateRange) {
			ResponseWithError(h.logger, w, err.Error(), http.StatusBadRequest)
			return
		}

		ResponseWithError(h.logger, w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("cost calculated",
		"total_cost", response.TotalCost,
		"start_date", response.StartDate,
		"end_date", response.EndDate)

	Response(h.logger, w, response, http.StatusOK)
}
