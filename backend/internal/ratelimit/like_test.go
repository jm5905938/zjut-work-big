package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type fakeScripter struct {
	result []int64
	err    error
}

func (f *fakeScripter) command(ctx context.Context) *redis.Cmd {
	cmd := redis.NewCmd(ctx)
	if f.err != nil {
		cmd.SetErr(f.err)
		return cmd
	}

	values := make([]any, len(f.result))
	for i, value := range f.result {
		values[i] = value
	}
	cmd.SetVal(values)

	return cmd
}

func (f *fakeScripter) Eval(
	ctx context.Context,
	_ string,
	_ []string,
	_ ...any,
) *redis.Cmd {
	return f.command(ctx)
}

func (f *fakeScripter) EvalSha(
	ctx context.Context,
	_ string,
	_ []string,
	_ ...any,
) *redis.Cmd {
	return f.command(ctx)
}

func (f *fakeScripter) EvalRO(
	ctx context.Context,
	_ string,
	_ []string,
	_ ...any,
) *redis.Cmd {
	return f.command(ctx)
}

func (f *fakeScripter) EvalShaRO(
	ctx context.Context,
	_ string,
	_ []string,
	_ ...any,
) *redis.Cmd {
	return f.command(ctx)
}

func (f *fakeScripter) ScriptExists(
	ctx context.Context,
	_ ...string,
) *redis.BoolSliceCmd {
	return redis.NewBoolSliceCmd(ctx)
}

func (f *fakeScripter) ScriptLoad(
	ctx context.Context,
	_ string,
) *redis.StringCmd {
	return redis.NewStringCmd(ctx)
}

func TestLikeLimiterAllow(t *testing.T) {
	tests := []struct {
		name        string
		result      []int64
		limit       int64
		wantAllowed bool
		wantRetry   time.Duration
		wantError   bool
	}{
		{
			name:        "within limit",
			result:      []int64{60, 1500},
			limit:       60,
			wantAllowed: true,
			wantRetry:   1500 * time.Millisecond,
		},
		{
			name:        "over limit",
			result:      []int64{61, 1200},
			limit:       60,
			wantAllowed: false,
			wantRetry:   1200 * time.Millisecond,
		},
		{
			name:      "invalid script result",
			result:    []int64{1},
			limit:     60,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewLikeLimiter(
				&fakeScripter{result: tt.result},
				tt.limit,
				time.Minute,
			)

			allowed, retryAfter, err := limiter.Allow(
				context.Background(),
				5,
			)
			if (err != nil) != tt.wantError {
				t.Fatalf("Allow() error = %v, wantError = %v", err, tt.wantError)
			}
			if allowed != tt.wantAllowed {
				t.Errorf("Allow() allowed = %v, want %v", allowed, tt.wantAllowed)
			}
			if retryAfter != tt.wantRetry {
				t.Errorf("Allow() retryAfter = %v, want %v", retryAfter, tt.wantRetry)
			}
		})
	}
}

func TestLikeLimiterAllowReturnsRedisError(t *testing.T) {
	limiter := NewLikeLimiter(
		&fakeScripter{err: errors.New("redis unavailable")},
		60,
		time.Minute,
	)

	_, _, err := limiter.Allow(context.Background(), 5)
	if err == nil {
		t.Fatal("Allow() error = nil, want non-nil")
	}
}
