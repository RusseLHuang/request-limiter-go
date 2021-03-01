package main

import (
	"github.com/RusseLHuang/request-limiter/limiter"
	"github.com/RusseLHuang/request-limiter/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	requestLimit := viper.GetInt("limit")
	requestLimitDuration := viper.GetInt("limitDuration")

	redisClient := redis.GetClient()
	redisRepository := limiter.NewRedisRepository(redisClient)
	limiterService := limiter.NewLimiterService(
		redisRepository,
		requestLimit,
		requestLimitDuration,
	)
	limiterController := limiter.NewLimiterController(limiterService)

	r := gin.Default()

	r.GET("/", limiterController.Limit)

	r.Run()
}
