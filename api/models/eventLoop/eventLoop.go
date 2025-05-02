package eventLoop

import (
	"os"
	"os/signal"
	"pledge-backend/api/ethereum"
	"pledge-backend/db"
	"pledge-backend/log"
	"syscall"
	"time"
)

const BlockLatest = "latest"
const BlockEarliest = "earliest"
const BlockPending = "pending"

var blockArgs = []string{BlockLatest, BlockEarliest, BlockPending}

// Dispatch 需求d  需求e（limit.go）
func Dispatch() {
	ticker := time.NewTicker(110 * time.Second)
	defer ticker.Stop()

	// 创建一个用于退出的channel
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			// 定时触发
			// 循环调用controllers.GetBlockByTag() 分别传入："latest", "earliest", "pending"
			for _, arg := range blockArgs {
				_, err := cacheBlockByTag(arg)
				if err != nil {
					log.Logger.Error("CacheBlockByTag error:" + err.Error())
					continue
				}
				log.Logger.Info("cache " + arg + " success")
			}
		case <-stopChan:
			log.Logger.Info("Received shutdown signal, exiting event loop...")
			return
		}
	}
}

func cacheBlockByTag(tag string) (*map[string]interface{}, error) {
	var arg string
	//把tag赋值给 arg，tag == "head"，arg = "latest"
	if tag == "head" {
		arg = BlockLatest
	} else {
		arg = tag
	}

	// redis里没有数据则从链上获取
	_blockMap, err := ethereum.GetBlockByRpcCall(arg)
	if err != nil {
		return nil, err
	}

	// 取到设置到Redis
	err = db.RedisSet(tag, _blockMap, 120)
	if err != nil {
		log.Logger.Error("Unmarshal json failed:" + err.Error())
	}
	return _blockMap, nil
}
