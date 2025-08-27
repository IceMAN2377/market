package service

import (
	"context"
	"github.com/IceMAN2377/market/internal/models"
)

type Service interface {
	CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error)
	GetSubscriptions(ctx context.Context, filters *models.SubscriptionFilters) (*models.SubscriptionListResponse, error)
	UpdateSubscription(ctx context.Context, id int, req *models.UpdateSubscriptionRequest) (*models.Subscription, error)
	DeleteSubscription(ctx context.Context, id int) error
	CalculateCost(ctx context.Context, req *models.CostCalculationRequest) (*models.CostCalculationResponse, error)
}
