package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"subscription-service/internal/app/entity"
	"subscription-service/internal/port"
)

const (
	maxRetries             = 3
	defaultInitialInterval = 3 * time.Millisecond
)

var _ SubscriptionUseCase = (*Subscription)(nil)

type Subscription struct {
	subscriptionRepo      port.SubscriptionRepo
	transactionController port.TransactionController
	backoffTxDoer         backoff.BackOff
	logger                *zap.Logger
}

func NewSubscription(
	subscriptionRepo port.SubscriptionRepo,
	transactionController port.TransactionController,
	logger *zap.Logger,
) (*Subscription, error) {
	expBackoffDoer := backoff.WithMaxRetries(backoff.NewExponentialBackOff(
		backoff.WithInitialInterval(defaultInitialInterval),
	), maxRetries)
	return &Subscription{
		subscriptionRepo:      subscriptionRepo,
		transactionController: transactionController,
		backoffTxDoer:         expBackoffDoer,
		logger:                logger,
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
		if errors.Is(err, port.ErrSubscriptionAlreadyExists) {
			return nil, ErrSubscriptionAlreadyExists
		}
		return nil, err
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
	bErr := backoff.Retry(
		func() error {
			err := r.update(ctx, post)
			if err != nil {
				if errors.Is(err, port.ErrSubscriptionAlreadyExists) {
					return ErrSubscriptionAlreadyExists
				}
				if !errors.Is(err, port.ErrTransactionFailure) {
					return backoff.Permanent(err)
				}

				return err
			}

			return nil
		},
		r.backoffTxDoer,
	)

	return bErr
}

func (r *Subscription) update(ctx context.Context, post entity.UpdateSubscriptionRequest) error {
	tx, err := r.transactionController.BeginTx(ctx, entity.RepeatableRead)
	if err != nil {
		return backoff.Permanent(err)
	}

	err = r.subscriptionRepo.Update(ctx, post)
	if err != nil {
		if errors.Is(err, port.ErrSubscriptionAlreadyExists) {
			return err
		}
		errR := tx.Rollback(ctx)
		if errR != nil {
			r.logger.Error("transaction rollback failed", zap.Error(err))
		}

		if errors.Is(err, port.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("update subscription: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("update trancastion: %w", err)
	}

	return nil
}

func (r *Subscription) Delete(ctx context.Context, id string) error {
	err := r.subscriptionRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			return ErrNotFound
		}
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
