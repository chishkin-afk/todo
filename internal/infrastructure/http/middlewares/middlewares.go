package middlewares

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/chishkin-afk/todo/internal/application/dtos"
	"github.com/chishkin-afk/todo/pkg/consts"
	errs "github.com/chishkin-afk/todo/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type jwtManager interface {
	Validate(token string) (uuid.UUID, error)
}

func AuthMiddleware(jm jwtManager, noAuth map[string]bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if noAuth[ctx.FullPath()] {
			ctx.Next()
			return
		}

		token := ctx.GetHeader("Authorization")
		userID, err := jm.Validate(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrMsg{
				Error: errs.ErrInvalidToken.Error(),
			})
			return
		}

		ctxValue := context.WithValue(ctx, consts.UserID, userID)
		ctx.Request = ctx.Request.WithContext(ctxValue)

		ctx.Next()
	}
}

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	times, exists := rl.requests[key]
	if !exists {
		times = []time.Time{}
	}

	var validTimes []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) >= rl.limit {
		rl.requests[key] = validTimes
		return false
	}

	validTimes = append(validTimes, now)
	rl.requests[key] = validTimes

	return true
}

func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	rl := newRateLimiter(limit, window)

	return func(ctx *gin.Context) {
		key := ctx.ClientIP()
		if uid, ok := ctx.Get(consts.UserID); ok {
			if id, ok := uid.(uuid.UUID); ok {
				key = id.String()
			}
		}

		if !rl.allow(key) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, dtos.ErrMsg{
				Error: errs.ErrTooManyRequests.Error(),
			})
			return
		}

		ctx.Next()
	}
}
