package errors_package

import "errors"

var (
	NoVariable         = errors.New("no such variable in env")
	EmptySub           = errors.New("subscription cannot be empty")
	EmptySubName       = errors.New("subscription name cannot be empty")
	InvalidPrice       = errors.New("price cannot be less than zero")
	SubscriptionInPast = errors.New("subscription in past")
	EmptyUser          = errors.New("user cannot be empty")
	ErrEmptyId         = errors.New("id cannot be empty")
	InvalidPeriod      = errors.New("invalid period")
	ErrDateRequired    = errors.New("date is required")
)
