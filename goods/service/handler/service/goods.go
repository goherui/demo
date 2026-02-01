package service

import (
	"context"
	"demo/goods/service/basic/config"
	__ "demo/goods/service/basic/proto"
	"demo/goods/service/model"
	"net/http"
	"time"

	"github.com/gospacex/gospacex/core/storage/cache/redis"
)

func (s *Server) GoodsCreate(_ context.Context, in *__.GoodsCreateReq) (*__.GoodsCreateResp, error) {
	Ctx := context.Background()
	Rdb := redis.RC
	result, err := Rdb.Get(Ctx, "goods"+in.Title).Result()
	if result != "" {
		return &__.GoodsCreateResp{
			Msg:   "商品已存在 ",
			Code:  http.StatusBadRequest,
			Goods: result,
		}, nil
	}
	Rdb.Set(Ctx, "goods"+in.Title, 1, time.Hour*2)
	goods := model.Goods{
		Title: in.Title,
		Price: float64(in.Price),
		Stock: int(in.Stock),
	}
	err = goods.GoodsCreate(config.DB)
	if err != nil {
		return nil, err
	}
	return &__.GoodsCreateResp{
		Msg:  "添加成功 ",
		Code: http.StatusBadRequest,
	}, nil
}
