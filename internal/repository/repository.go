package repository

import (
	"context"
	"github.com/IceMAN2377/market/internal/models"
)

type Repository interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error)
	GetSubscriptions(ctx context.Context, filters *models.SubscriptionFilters) ([]models.Subscription, error)
	UpdateSubscription(ctx context.Context, id int, updates *models.UpdateSubscriptionRequest) (*models.Subscription, error)
	DeleteSubscription(ctx context.Context, id int) error
	CalculateCost(ctx context.Context, req *models.CostCalculationRequest) (int, error)
}
