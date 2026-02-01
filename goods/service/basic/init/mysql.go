package init

import (
	"demo/goods/service/basic/config"
	"demo/goods/service/basic/logger"
	"demo/goods/service/model"
	"time"

	"github.com/gospacex/gospacex/core/storage/conf"
	"github.com/gospacex/gospacex/core/storage/db/mysql"
	"go.uber.org/zap"
)

func InitMySQL() {
	var err error
	config.DB, err = mysql.Init(false, "debug", conf.Cfg.Mysql)
	if err != nil {
		logger.Error("mysql连接失败", zap.Error(err))
		return
	}
	logger.Info("mysql连接成功")
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, _ := config.DB.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	err = config.DB.AutoMigrate(&model.User{}, model.Goods{})
	if err != nil {
		logger.Error("迁移失败", zap.Error(err))
		return
	}
	logger.Info("迁移成功")

}
