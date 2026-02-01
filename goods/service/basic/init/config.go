package init

import (
	"demo/goods/service/basic/logger"

	"github.com/gospacex/gospacex/core/storage/conf"
)

func InitConfig() {
	conf.ParseConfig("../../../")
	logger.Info("配置加载成功")
}
