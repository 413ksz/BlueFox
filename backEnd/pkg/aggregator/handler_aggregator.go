package aggregator

import "github.com/413ksz/BlueFox/backEnd/domain/users"

// HandlerAggregator aggregates handlers from different packages into a single struct for easy access for the router.
type HandlerAggregator struct {
	UserHandler users.UserHandler
}

// NewHandlerAggregator creates a new HandlerAggregator with the provided handlers.
func NewHandlerAggregator(userHandler users.UserHandler) *HandlerAggregator {
	return &HandlerAggregator{
		UserHandler: userHandler,
	}
}
