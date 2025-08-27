package subscription

import (
	"context"
	errs "github.com/IceMAN2377/market/internal/errors"
	"github.com/IceMAN2377/market/internal/models"
	"github.com/IceMAN2377/market/internal/repository"
	"github.com/IceMAN2377/market/internal/service"
	"github.com/google/uuid"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type subscription struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) service.Service {
	return &subscription{
		repo: repo,
	}
}

// validateDateFormat проверяет формат даты MM-YYYY
func (s *subscription) validateDateFormat(date string) error {
	matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-\d{4}$`, date)
	if !matched {
		return errs.ErrInvalidDateFormat
	}

	// Дополнительная проверка - парсим дату
	parts := strings.Split(date, "-")
	month, _ := strconv.Atoi(parts[0])
	year, _ := strconv.Atoi(parts[1])

	if month < 1 || month > 12 {
		return errs.ErrInvalidDateFormat
	}

	currentYear := time.Now().Year()
	if year < 2000 || year > currentYear+10 {
		return errs.ErrInvalidDateFormat
	}

	return nil
}

// validateUUID проверяет корректность UUID
func (s *subscription) validateUUID(uuidStr string) error {
	if _, err := uuid.Parse(uuidStr); err != nil {
		return errs.ErrInvalidUUID
	}
	return nil
}

// validateDateRange проверяет корректность диапазона дат
func (s *subscription) validateDateRange(startDate, endDate string) error {
	if err := s.validateDateFormat(startDate); err != nil {
		return err
	}

	if err := s.validateDateFormat(endDate); err != nil {
		return err
	}

	// Сравниваем даты (простое строковое сравнение работает для формата MM-YYYY)
	startParts := strings.Split(startDate, "-")
	endParts := strings.Split(endDate, "-")

	startYear, _ := strconv.Atoi(startParts[1])
	startMonth, _ := strconv.Atoi(startParts[0])
	endYear, _ := strconv.Atoi(endParts[1])
	endMonth, _ := strconv.Atoi(endParts[0])

	startTime := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(endYear, time.Month(endMonth), 1, 0, 0, 0, 0, time.UTC)

	if startTime.After(endTime) {
		return errs.ErrInvalidDateRange
	}

	return nil
}

func (s *subscription) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	// Валидация UUID пользователя
	if err := s.validateUUID(req.UserID); err != nil {
		return nil, err
	}

	// Валидация формата начальной даты
	if err := s.validateDateFormat(req.StartDate); err != nil {
		return nil, err
	}

	// Валидация конечной даты, если она указана
	if req.EndDate != nil {
		if err := s.validateDateRange(req.StartDate, *req.EndDate); err != nil {
			return nil, err
		}
	}

	// Валидация цены
	if req.Price <= 0 {
		return nil, errs.ErrInvalidPrice
	}

	// Создание модели подписки
	subscription := &models.Subscription{
		ServiceName: strings.TrimSpace(req.ServiceName),
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	return s.repo.CreateSubscription(ctx, subscription)
}

func (s *subscription) GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error) {
	if id <= 0 {
		return nil, errs.ErrInvalidData
	}

	return s.repo.GetSubscriptionByID(ctx, id)
}

func (s *subscription) GetSubscriptions(ctx context.Context, filters *models.SubscriptionFilters) (*models.SubscriptionListResponse, error) {
	// Валидация параметров пагинации
	if filters.Limit <= 0 {
		filters.Limit = 10
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	// Валидация UUID пользователя, если указан
	if filters.UserID != nil {
		if err := s.validateUUID(*filters.UserID); err != nil {
			return nil, err
		}
	}

	// Получение подписок
	subscriptions, err := s.repo.GetSubscriptions(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Для простоты не считаем total (можно добавить отдельный метод в репозиторий)
	response := &models.SubscriptionListResponse{
		Subscriptions: subscriptions,
		Total:         len(subscriptions), // Упрощение
		Limit:         filters.Limit,
		Offset:        filters.Offset,
	}

	return response, nil
}

func (s *subscription) UpdateSubscription(ctx context.Context, id int, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	if id <= 0 {
		return nil, errs.ErrInvalidData
	}

	// Проверяем, что подписка существует
	existing, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Валидация конечной даты, если она обновляется
	if req.EndDate != nil {
		if err := s.validateDateRange(existing.StartDate, *req.EndDate); err != nil {
			return nil, err
		}
	}

	// Валидация цены, если она обновляется
	if req.Price != nil && *req.Price <= 0 {
		return nil, errs.ErrInvalidPrice
	}

	// Обрезаем пробелы в названии сервиса
	if req.ServiceName != nil {
		trimmed := strings.TrimSpace(*req.ServiceName)
		req.ServiceName = &trimmed
	}

	return s.repo.UpdateSubscription(ctx, id, req)
}

func (s *subscription) DeleteSubscription(ctx context.Context, id int) error {
	if id <= 0 {
		return errs.ErrInvalidData
	}

	return s.repo.DeleteSubscription(ctx, id)
}

func (s *subscription) CalculateCost(ctx context.Context, req *models.CostCalculationRequest) (*models.CostCalculationResponse, error) {
	// Валидация диапазона дат
	if err := s.validateDateRange(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	// Валидация UUID пользователя, если указан
	if req.UserID != nil {
		if err := s.validateUUID(*req.UserID); err != nil {
			return nil, err
		}
	}

	// Расчет стоимости
	totalCost, err := s.repo.CalculateCost(ctx, req)
	if err != nil {
		return nil, err
	}

	response := &models.CostCalculationResponse{
		TotalCost:   totalCost,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
	}

	return response, nil
}
