package usecase

import (
	"context"

	"subscription-service/internal/app/entity"
)

type SubscriptionUseCase interface {
	Create(ctx context.Context, post entity.CreateSubscriptionRequest) (*entity.Subscription, error)
	Read(ctx context.Context, id string) (*entity.Subscription, error)
	Update(ctx context.Context, post entity.UpdateSubscriptionRequest) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter entity.ListSubscriptionFilter) ([]entity.Subscription, error)
	Sum(ctx context.Context, filter entity.ListSubscriptionFilter) (int64, error)
}
