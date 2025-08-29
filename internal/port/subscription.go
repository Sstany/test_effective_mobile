package port

import (
	"context"

	"subscription-service/internal/app/entity"
)

//go:generate mockgen -destination ../adapter/repo/mock/subscription_mock.go -package repo -source ./subscription.go

type SubscriptionRepo interface {
	Create(ctx context.Context, post entity.CreateSubscriptionRequest) error
	GetSubscription(ctx context.Context, id string) (*entity.Subscription, error)
	Update(ctx context.Context, post entity.UpdateSubscriptionRequest) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter entity.ListSubscriptionFilter) ([]entity.Subscription, error)
	Sum(ctx context.Context, filter entity.ListSubscriptionFilter) (int64, error)
}
