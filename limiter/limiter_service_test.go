package limiter

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type LimiterRepositoryMock struct {
	RequestCount map[string]interface{}
	Expiration   map[string]int64
}

func (r *LimiterRepositoryMock) Increment(ctx context.Context, key string) int {
	if r.RequestCount[key] == nil {
		r.RequestCount[key] = 0
	}

	r.RequestCount[key] = r.RequestCount[key].(int) + 1
	return r.RequestCount[key].(int)
}

func (r *LimiterRepositoryMock) Get(ctx context.Context, key string) string {
	now := time.Now().Unix()

	if r.Expiration[key] == 0 {
		if r.RequestCount[key] == nil {
			return ""
		}
		val := fmt.Sprintf("%v", r.RequestCount[key])
		return val
	}

	if r.Expiration[key] < now || r.RequestCount[key] == nil {
		return ""
	}

	return fmt.Sprintf("%v", r.RequestCount[key])
}

func (r *LimiterRepositoryMock) Set(ctx context.Context, key string, value string, duration int) {
	timeDuration := time.Duration(duration) * time.Duration(time.Second)
	r.RequestCount[key] = value
	r.Expiration[key] = time.Now().Add(timeDuration).Unix()
}

func (r *LimiterRepositoryMock) Del(ctx context.Context, key string) {
	r.RequestCount[key] = 0
}

func TestRequestCountWithinDurationAndLimit(t *testing.T) {
	repository := &LimiterRepositoryMock{
		RequestCount: make(map[string]interface{}),
		Expiration:   make(map[string]int64),
	}
	requestLimiter := NewLimiterService(
		repository,
		5,
		5,
	)

	ctx := context.Background()
	sourceIP := "192.168.0.1"

	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)
	resp, err := requestLimiter.LimitRequest(ctx, sourceIP)
	if err != nil {
		t.Fatalf("Error should not happened when limit not exceeded yet")
	}

	if resp != 4 {
		t.Fatalf("Request count is not correct %d", resp)
	}
}

func TestRequestCountExceedLimit(t *testing.T) {
	repository := &LimiterRepositoryMock{
		RequestCount: make(map[string]interface{}),
		Expiration:   make(map[string]int64),
	}
	requestLimiter := NewLimiterService(
		repository,
		5,
		5,
	)

	ctx := context.Background()
	sourceIP := "192.168.0.1"

	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)
	resp, err := requestLimiter.LimitRequest(ctx, sourceIP)
	if err != nil {
		t.Fatalf("Error should not happened when limit not exceeded yet")
	}

	if resp != 5 {
		t.Fatalf("Request count is not correct %d", resp)
	}

	resp, err = requestLimiter.LimitRequest(ctx, sourceIP)
	if err.Error() != "Error" {
		t.Fatalf("Should return error message when exceed request limit")
	}
}

func TestRequestCountShouldRefreshAfterDuration(t *testing.T) {
	repository := &LimiterRepositoryMock{
		RequestCount: make(map[string]interface{}),
		Expiration:   make(map[string]int64),
	}
	duration := 2
	durationAfterTimeout := duration + 1
	requestLimiter := NewLimiterService(
		repository,
		5,
		duration,
	)

	ctx := context.Background()
	sourceIP := "192.168.0.1"

	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)
	requestLimiter.LimitRequest(ctx, sourceIP)

	time.Sleep(time.Duration(durationAfterTimeout) * time.Duration(time.Second))

	resp, err := requestLimiter.LimitRequest(ctx, sourceIP)
	if err != nil {
		t.Fatalf("Error should not happened when limit not exceeded yet")
	}

	if resp != 1 {
		t.Fatalf("Request count is not correct %d", resp)
	}
}
