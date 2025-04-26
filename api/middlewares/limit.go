package middlewares

import (
	"github.com/gin-gonic/gin"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/response"
	"time"
)

// 需求e

// 存储IP和其请求的时间戳
var ipRequests = make(map[string][]time.Time)

// 每秒最大请求数
const maxRequests = 1
const timeWindow = time.Second // 时间窗口，1秒钟

func Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := response.Gin{Res: c}
		// 获取客户端IP
		ip := c.ClientIP()

		// 获取当前时间戳
		now := time.Now()

		// 从存储中获取该IP的所有请求时间戳
		timestamps, exists := ipRequests[ip]
		if !exists {
			timestamps = []time.Time{}
		}

		// 移除过期的时间戳（超过1秒的请求）
		var validTimestamps []time.Time
		for _, timestamp := range timestamps {
			if now.Sub(timestamp) < timeWindow {
				validTimestamps = append(validTimestamps, timestamp)
			}
		}

		// 更新该IP的请求时间戳
		ipRequests[ip] = append(validTimestamps, now)

		// 判断该IP是否超过最大请求次数
		if len(ipRequests[ip]) > maxRequests {
			// 超过限制，返回限流错误
			res.Response(c, statecode.RequestLimit, nil)
			// 终止请求
			c.Abort()
			return
		}

		// 允许请求，继续执行
		c.Next()
	}
}
