package entity

import "time"

type Subscription struct {
	ID        string
	Title     string
	Price     int64
	UserID    string
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt int64
	UpdatedAt int64
}

type UpdateSubscriptionRequest struct {
	ID        string
	Title     string
	Price     int64
	UserID    string
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt int64
	UpdatedAt int64
}

type CreateSubscriptionRequest UpdateSubscriptionRequest

type ListSubscriptionFilter struct {
	Title     *string
	UserID    *string
	Price     *int64
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int
	Offset    *int
}
