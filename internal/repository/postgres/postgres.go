package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	errs "github.com/IceMAN2377/market/internal/errors"
	"github.com/IceMAN2377/market/internal/models"
	"github.com/IceMAN2377/market/internal/repository"
	"github.com/jmoiron/sqlx"
)

type postgres struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.Repository {
	return &postgres{
		db: db,
	}
}

func (p *postgres) CreateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`

	var result models.Subscription
	err := p.db.GetContext(ctx, &result,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return &result, nil
}

func (p *postgres) GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1`

	var subscription models.Subscription
	err := p.db.GetContext(ctx, &subscription, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &subscription, nil
}

func (p *postgres) GetSubscriptions(ctx context.Context, filters *models.SubscriptionFilters) ([]models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filters.UserID)
		argIndex++
	}

	if filters.ServiceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.ServiceName+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, filters.Offset)

	var subscriptions []models.Subscription
	err := p.db.SelectContext(ctx, &subscriptions, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if subscriptions == nil {
		subscriptions = []models.Subscription{}
	}

	return subscriptions, nil
}

func (p *postgres) UpdateSubscription(ctx context.Context, id int, updates *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	var setParts []string
	var args []interface{}
	argIndex := 1

	if updates.ServiceName != nil {
		setParts = append(setParts, fmt.Sprintf("service_name = $%d", argIndex))
		args = append(args, *updates.ServiceName)
		argIndex++
	}

	if updates.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *updates.Price)
		argIndex++
	}

	if updates.EndDate != nil {
		setParts = append(setParts, fmt.Sprintf("end_date = $%d", argIndex))
		args = append(args, *updates.EndDate)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, errs.ErrInvalidData
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	query := fmt.Sprintf(`
		UPDATE subscriptions 
		SET %s 
		WHERE id = $%d
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	args = append(args, id)

	var subscription models.Subscription
	err := p.db.GetContext(ctx, &subscription, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return &subscription, nil
}

func (p *postgres) DeleteSubscription(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (p *postgres) CalculateCost(ctx context.Context, req *models.CostCalculationRequest) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0) as total_cost
		FROM subscriptions
		WHERE 1=1`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if req.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.ServiceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name ILIKE $%d", argIndex))
		args = append(args, "%"+*req.ServiceName+"%")
		argIndex++
	}

	// Добавляем условия для периода
	// Подписка пересекается с запрашиваемым периодом если:
	// start_date <= end_period AND (end_date IS NULL OR end_date >= start_period)
	conditions = append(conditions, fmt.Sprintf("start_date <= $%d", argIndex))
	args = append(args, req.EndDate)
	argIndex++

	conditions = append(conditions, fmt.Sprintf("(end_date IS NULL OR end_date >= $%d)", argIndex))
	args = append(args, req.StartDate)
	argIndex++

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var totalCost int
	err := p.db.GetContext(ctx, &totalCost, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate cost: %w", err)
	}

	return totalCost, nil
}
