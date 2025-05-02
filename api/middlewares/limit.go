package middlewares

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/response"
	"pledge-backend/log"
	"sync"
	"syscall"
	"time"
)

// 需求e 限流每秒2个请求

// 便于测试，把桶的容量设置为1
var rateLimit = 1

// 每500毫秒放一个令牌
var tokenInterval = 500 * time.Millisecond

// 创建一个容量为rateLimit的桶
var (
	bucket = make(chan struct{}, rateLimit)
	once   sync.Once
)

// 初始化填满令牌
func initToken() {
	// 初始填充令牌
	for i := 0; i < rateLimit; i++ {
		bucket <- struct{}{}
		log.Logger.Info("初始化桶")
	}

	// 创建一个用于退出的channel
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// 每tokenInterval时间放一个令牌
		ticker := time.NewTicker(tokenInterval)
		defer ticker.Stop()

		for range ticker.C {
			select {
			case bucket <- struct{}{}:
				// 放入成功
				log.Logger.Info("放入成功")
			case <-stopChan:
				// 收到退出信号，退出循环
				return
			default:
				// 桶满了，不放
			}
		}
	}()
}

func Limit() gin.HandlerFunc {

	once.Do(initToken)
	return func(c *gin.Context) {
		res := response.Gin{Res: c}

		select {
		case <-bucket:
			// 获取令牌，允许请求
			break
		default:
			// 桶为空，返回限流错误
			res.Response(c, statecode.RequestLimit, nil)
			// 终止请求
			c.Abort()
			return
		}

		// 允许请求，继续执行
		c.Next()
	}
}
