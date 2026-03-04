package wx_official

import (
	wx_api "crontab/pkg/wx-api"
	"fmt"
)

type wxOfficial struct {
	*wx_api.DefaultToken
}

func NewWxOfficial(appid, secret string) wx_api.Token {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, secret)
	return &wxOfficial{
		DefaultToken: &wx_api.DefaultToken{
			AppId:  appid,
			Secret: secret,
			Url:    url,
		},
	}
}
