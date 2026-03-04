package cron

import (
	"crontab/pkg/config"
	"crontab/pkg/log"
	wx_official "crontab/pkg/wx-api/wx-official"
	"github.com/robfig/cron/v3"
)

var tables = []string{"url_map", "url_map_user"}

func Run() {
	cnf := config.GetConfig()
	GetWxAccessToken(cnf)()
	c := cron.New()
	c.AddFunc("*/5 * * * *", GetWxAccessToken(cnf))
	c.Run()
}

func GetWxAccessToken(cnf *config.Config) func() {
	return func() {
		for _, item := range cnf.WxOfficials {
			official := wx_official.NewWxOfficial(item.AppId, item.Secret)
			err := official.RefreshToken()
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}
