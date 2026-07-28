package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type fakeLikeRateLimiter struct {
	allowed    bool
	retryAfter time.Duration
	err        error
	called     bool
	userID     uint64
}

func (f *fakeLikeRateLimiter) Allow(
	_ context.Context,
	userID uint64,
) (bool, time.Duration, error) {
	f.called = true
	f.userID = userID

	return f.allowed, f.retryAfter, f.err
}

func TestLikeRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setUser        bool
		limiter        *fakeLikeRateLimiter
		wantStatus     int
		wantRetryAfter string
		wantHandler    bool
		wantLimiter    bool
	}{
		{
			name:        "allowed",
			setUser:     true,
			limiter:     &fakeLikeRateLimiter{allowed: true},
			wantStatus:  http.StatusNoContent,
			wantHandler: true,
			wantLimiter: true,
		},
		{
			name:           "rate limited",
			setUser:        true,
			limiter:        &fakeLikeRateLimiter{retryAfter: 1500 * time.Millisecond},
			wantStatus:     http.StatusTooManyRequests,
			wantRetryAfter: "2",
			wantLimiter:    true,
		},
		{
			name:        "redis unavailable",
			setUser:     true,
			limiter:     &fakeLikeRateLimiter{err: errors.New("redis unavailable")},
			wantStatus:  http.StatusServiceUnavailable,
			wantLimiter: true,
		},
		{
			name:       "missing user",
			limiter:    &fakeLikeRateLimiter{allowed: true},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(ErrorHandler())

			handlerCalled := false
			r.POST(
				"/posts/:post_id/like",
				func(c *gin.Context) {
					if tt.setUser {
						c.Set(ContextUserID, uint64(5))
					}
					c.Next()
				},
				LikeRateLimit(tt.limiter),
				func(c *gin.Context) {
					handlerCalled = true
					c.Status(http.StatusNoContent)
				},
			)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				http.MethodPost,
				"/posts/1/like",
				nil,
			)
			r.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if recorder.Header().Get("Retry-After") != tt.wantRetryAfter {
				t.Errorf(
					"Retry-After = %q, want %q",
					recorder.Header().Get("Retry-After"),
					tt.wantRetryAfter,
				)
			}
			if handlerCalled != tt.wantHandler {
				t.Errorf("handlerCalled = %v, want %v", handlerCalled, tt.wantHandler)
			}
			if tt.limiter.called != tt.wantLimiter {
				t.Errorf("limiter.called = %v, want %v", tt.limiter.called, tt.wantLimiter)
			}
			if tt.wantLimiter && tt.limiter.userID != 5 {
				t.Errorf("limiter.userID = %d, want 5", tt.limiter.userID)
			}
		})
	}
}
