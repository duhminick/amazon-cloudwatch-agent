package pusher

import "golang.org/x/time/rate"

type rateLimit struct {
	limiter rate.Limiter
}

func NewRateLimit(limit float64) *rateLimit {
	return &rateLimit{
		limiter: *rate.NewLimiter(1, 1),
	}
}
