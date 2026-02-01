package config

import (
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	Es *elastic.Client
)
