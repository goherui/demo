package init

import "demo/goods/service/basic/logger"

func init() {
	logger.InitLogger()
	InitConfig()
	InitMySQL()
	InitRedis()
	InitEs()
}
