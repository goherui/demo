package init

import (
	"demo/goods/service/basic/logger"

	"github.com/gospacex/gospacex/core/storage/cache/redis"
	"github.com/gospacex/gospacex/core/storage/conf"
	"go.uber.org/zap"
)

func InitRedis() {
	err := redis.Init(false, conf.Cfg.Redis)
	if err != nil {
		logger.Error("redis连接失败", zap.Error(err))
		return
	}
	logger.Info("redis连接成功")
}
