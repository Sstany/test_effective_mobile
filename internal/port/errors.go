package port

import "errors"

var (
	ErrNotFound                  = errors.New("subscription not found")
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")

	ErrTransactionFailure = errors.New("transaction failure")
)
