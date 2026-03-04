package server

import (
	"context"
	"crontab/pkg/config"
	"crontab/pkg/log"
	wx_official "crontab/pkg/wx-api/wx-official"
	"crontab/pkg/zerror"
	"crontab/proto"
)

type tokenServer struct {
	proto.UnimplementedTokenServer
	config *config.Config
	log    log.ILogger
}

func NewTokenServer(config *config.Config, log log.ILogger) proto.TokenServer {
	return &tokenServer{
		config: config,
		log:    log,
	}
}
func (s *tokenServer) GetToken(ctx context.Context, in *proto.TokenRequest) (*proto.TokenResponse, error) {
	secret := s.getSecret(in)
	if secret == "" {
		err := zerror.NewByMsg("令牌获取失败")
		s.log.Error(err)
		return nil, err
	}
	token := wx_official.NewWxOfficial(in.Appid, secret)
	accessToken, err := token.GetToken()
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	res := &proto.TokenResponse{
		AccessToken: accessToken.AccessToken,
	}
	return res, err
}
func (s *tokenServer) getSecret(in *proto.TokenRequest) string {
	for _, item := range s.config.WxOfficials {
		if item.AppId == in.Appid {
			return item.Secret
		}
	}
	return ""
}
