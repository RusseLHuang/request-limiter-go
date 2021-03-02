package main

import (
	"github.com/RusseLHuang/request-limiter/limiter"
	"github.com/RusseLHuang/request-limiter/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init() *gin.Engine {
	viper.SetConfigType("env")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	requestLimit := viper.GetInt("LIMIT")
	requestLimitDuration := viper.GetInt("LIMIT_DURATION")
	redisUrl := viper.GetString("REDIS_ENDPOINT")
	redisPassword := viper.GetString("REDIS_PASSWORD")

	redisClient := redis.NewClient(redisUrl, redisPassword)
	redisRepository := limiter.NewRedisRepository(redisClient)
	limiterService := limiter.NewLimiterService(
		redisRepository,
		requestLimit,
		requestLimitDuration,
	)
	limiterController := limiter.NewLimiterController(limiterService)

	r := gin.Default()

	r.GET("/", limiterController.Limit)

	return r
}

func main() {
	r := Init()

	r.Run()
}
