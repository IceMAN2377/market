package models

import (
	"time"
)

// Subscription представляет основную модель подписки
type Subscription struct {
	ID          int       `json:"id" db:"id"`
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int       `json:"price" db:"price"`
	UserID      string    `json:"user_id" db:"user_id"`
	StartDate   string    `json:"start_date" db:"start_date"`
	EndDate     *string   `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSubscriptionRequest для создания подписки
type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" validate:"required,max=255"`
	Price       int     `json:"price" validate:"required,min=1"`
	UserID      string  `json:"user_id" validate:"required,uuid"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

// UpdateSubscriptionRequest для обновления подписки
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,max=255"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,min=1"`
	EndDate     *string `json:"end_date,omitempty"`
}

// SubscriptionFilters для фильтрации при получении списка подписок
type SubscriptionFilters struct {
	UserID      *string `json:"user_id,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
	Limit       int     `json:"limit" validate:"min=1,max=100"`
	Offset      int     `json:"offset" validate:"min=0"`
}

// CostCalculationRequest для расчета стоимости подписок за период
type CostCalculationRequest struct {
	UserID      *string `json:"user_id,omitempty" validate:"omitempty,uuid"`
	ServiceName *string `json:"service_name,omitempty"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     string  `json:"end_date" validate:"required"`
}

// CostCalculationResponse ответ на запрос расчета стоимости
type CostCalculationResponse struct {
	TotalCost   int     `json:"total_cost"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	UserID      *string `json:"user_id,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
}

// SubscriptionListResponse для ответа со списком подписок
type SubscriptionListResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Total         int            `json:"total"`
	Limit         int            `json:"limit"`
	Offset        int            `json:"offset"`
}
