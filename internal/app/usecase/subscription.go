package usecase

import (
	"context"
	"fmt"
	"subscrioption-service/internal/app/entity"
	"subscrioption-service/internal/port"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var _ SubscriptionUseCase = (*Subscription)(nil)

type Subscription struct {
	subscriptionRepo port.SubscriptionRepo
	pool             *pgxpool.Pool
	logger           *zap.Logger
}

func NewSubscription(
	subscriptionRepo port.SubscriptionRepo,
	pool *pgxpool.Pool,
	logger *zap.Logger,
) (*Subscription, error) {
	return &Subscription{
		subscriptionRepo: subscriptionRepo,
		pool:             pool,
		logger:           logger,
	}, nil
}

func (r *Subscription) Create(
	ctx context.Context,
	post entity.CreateSubscriptionRequest,
) (*entity.Subscription, error) {
	id := uuid.NewString()
	post.ID = id
	err := r.subscriptionRepo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return &entity.Subscription{
		ID:        id,
		Title:     post.Title,
		Price:     post.Price,
		UserID:    post.UserID,
		StartDate: post.StartDate,
		EndDate:   post.EndDate,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (r *Subscription) Read(ctx context.Context, id string) (*entity.Subscription, error) {
	sub, err := r.subscriptionRepo.GetSubscription(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub, nil
}

func (r *Subscription) Update(ctx context.Context, post entity.UpdateSubscriptionRequest) error {
	err := r.subscriptionRepo.Update(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (r *Subscription) Delete(ctx context.Context, id string) error {
	err := r.subscriptionRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

func (r *Subscription) List(ctx context.Context, filter entity.ListSubscriptionFilter) ([]entity.Subscription, error) {
	subs, err := r.subscriptionRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return subs, nil
}

func (r *Subscription) Sum(ctx context.Context, filter entity.ListSubscriptionFilter) (int64, error) {
	sum, err := r.subscriptionRepo.Sum(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to sum subscriptions: %w", err)
	}

	return sum, nil
}
