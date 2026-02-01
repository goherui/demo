package model

import "gorm.io/gorm"

type Goods struct {
	gorm.Model
	Title string  `gorm:"type:varchar(30);comment:商品标题"`
	Price float64 `gorm:"type:decimal(10,2);comment:商品价格"`
	Stock int     `gorm:"type:int;comment:商品标题"`
}

func (g *Goods) GoodsCreate(db *gorm.DB) error {
	return db.Create(&g).Error
}
