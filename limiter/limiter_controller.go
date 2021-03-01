package limiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LimiterController struct {
	LimiterService *LimiterService
}

func NewLimiterController(limiterService *LimiterService) *LimiterController {
	return &LimiterController{
		LimiterService: limiterService,
	}
}

func (limiter *LimiterController) Limit(c *gin.Context) {
	clientIP := c.ClientIP()
	response, err := limiter.LimiterService.LimitRequest(c, clientIP)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, response)
	}
}
