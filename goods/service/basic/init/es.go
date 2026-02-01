package init

import (
	"demo/goods/service/basic/config"
	"demo/goods/service/basic/logger"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func InitEs() {
	var err error
	config.Es, err = elastic.NewClient(elastic.SetURL("http://115.190.54.31:9200"),
		elastic.SetSniff(false))
	if err != nil {
		// Handle error
		logger.Error("es连接失败", zap.Error(err))
		return
	}
	logger.Info("es连接成功")
}
