package service

import (
	"context"
	"demo/goods/bff/basic/middleware"
	"demo/goods/service/basic/config"
	__ "demo/goods/service/basic/proto"
	"demo/goods/service/model"
	"net/http"
	"strconv"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	__.UnimplementedStreamGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) Login(_ context.Context, in *__.LoginReq) (*__.LoginResp, error) {
	var user model.User
	err := user.FindUser(config.DB, in.Username)
	if err != nil {
		return nil, err
	}
	if user.Password != in.Password {
		return &__.LoginResp{
			Msg:  "登录失败",
			Code: http.StatusBadRequest,
		}, nil
	}
	handler, err := middleware.TokenHandler(int(user.ID))
	if err != nil {
		return nil, err
	}
	userMap := map[string]string{
		"userId": strconv.Itoa(int(user.ID)),
		"token":  handler,
	}
	return &__.LoginResp{
		Msg:     "登录成功 ",
		Code:    http.StatusBadRequest,
		UserMap: userMap,
	}, nil
}
