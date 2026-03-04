package token_server

import (
	"crontab/pkg/config"
	"crontab/pkg/log"
	"crontab/proto"
	"crontab/token-server/server"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func Start(cnf *config.Config, logger log.ILogger) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port))
	if err != nil {
		log.Error(err)
		return
	}
	s := grpc.NewServer(server.GetOptions()...)
	tokenServer := server.NewTokenServer(cnf, logger)
	proto.RegisterTokenServer(s, tokenServer)

	healthCheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthCheck)

	if err := s.Serve(lis); err != nil {
		log.Error(err)
		return
	}

}
