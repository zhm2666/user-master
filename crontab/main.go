package main

import (
	"context"
	"crontab/cron"
	"crontab/pkg/config"
	"crontab/pkg/db/redis"
	"crontab/pkg/log"
	token_server "crontab/token-server"
	"flag"
	"os"
	"os/signal"
)

var (
	configFile = flag.String("config", "dev.config.yaml", "")
)

func main() {
	flag.Parse()
	//初始化配置文件
	config.InitConfig(*configFile)
	cnf := config.GetConfig()

	logger := log.NewLogger()

	logger.SetLevel(cnf.Log.Level)
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetPrintCaller(true)

	redis.InitRedisPool(cnf)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go cron.Run()
	go token_server.Start(cnf, logger)

	<-ctx.Done()
}
