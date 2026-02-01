package request

type GoodsCreate struct {
	Title string  `form:"title" json:"title" xml:"title" binding:"required"`
	Price float64 `form:"price" json:"price" xml:"price" binding:"required"`
	Stock int     `form:"stock" json:"stock" xml:"stock" binding:"required"`
}
