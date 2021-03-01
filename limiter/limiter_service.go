package limiter

import (
	"context"
	"errors"
	"strconv"
	"time"
)

type LimiterService struct {
	Repository           LimiterRepository
	RequestLimit         int
	RequestLimitDuration int
}

func NewLimiterService(
	repository LimiterRepository,
	requestLimit int,
	requestLimitDuration int,
) *LimiterService {
	return &LimiterService{
		Repository:           repository,
		RequestLimit:         requestLimit,
		RequestLimitDuration: requestLimitDuration,
	}
}

func (limiter *LimiterService) LimitRequest(
	ctx context.Context,
	sourceIP string,
) (int, error) {
	requestKey := sourceIP
	now := int(time.Now().Unix())
	requestLatestTimeKey := requestKey + ":time"
	requestLatestTime := limiter.Repository.Get(ctx, requestLatestTimeKey)

	if requestLatestTime == "" {
		limiter.Repository.Del(ctx, requestKey)
		limiter.Repository.Set(ctx, requestLatestTimeKey, strconv.Itoa(now), limiter.RequestLimitDuration)
	}

	limiter.Repository.Increment(ctx, requestKey)
	requestCount := limiter.Repository.Get(ctx, requestKey)
	requestCountInt, _ := strconv.Atoi(requestCount)

	if limiter.isExceeded(requestCountInt) {
		err := errors.New("Error")
		return 0, err
	}

	return requestCountInt, nil
}

func (limiter *LimiterService) isExceeded(requestCount int) bool {
	return requestCount > limiter.RequestLimit
}
