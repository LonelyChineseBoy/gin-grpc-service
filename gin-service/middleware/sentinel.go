package middleware

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func SentinelWarmUp() gin.HandlerFunc {
	return func(context *gin.Context) {
		_, err := flow.LoadRules([]*flow.Rule{
			{
				Resource:               "test",
				TokenCalculateStrategy: flow.WarmUp,
				ControlBehavior:        flow.Reject,
				Threshold:              10,
				WarmUpPeriodSec:        10,
				StatIntervalInMs:       3000,
			},
		})
		if err != nil {
			zap.S().Errorf("new flow rules failed:%s", err)
		}
		entry, blockError := sentinel.Entry("test", sentinel.WithTrafficType(base.Inbound), sentinel.WithResourceType(base.ResTypeWeb))
		if blockError != nil {
			context.JSON(http.StatusRequestTimeout, gin.H{
				"message": "限流了",
			})
		} else {
			context.Next()
		}
		defer entry.Exit()
	}
}
