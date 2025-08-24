package entity

import "time"

type Subscription struct {
	ID        string
	Title     string
	Price     int64
	UserID    string
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateSubscriptionRequest struct {
	Title     string
	Price     int64
	UserID    string
	StartDate time.Time
	EndDate   *time.Time
}

type CreateSubscriptionRequest UpdateSubscriptionRequest

type ListSubscriptionFilter struct {
	Title     *string
	UserID    *string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int64
	Offset    *int64
}
