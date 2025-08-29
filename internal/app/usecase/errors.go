package usecase

import "errors"

var (
	ErrSubscriptionNotFound      = errors.New("subscription not found")
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	ErrInvalidSubscriptionData   = errors.New("invalid subscription data")

	ErrNotFound           = errors.New("subscription not found")
	ErrTransactionFailure = errors.New("transaction failure")
)
