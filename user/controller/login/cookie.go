package login

import (
	"net/http"
	"time"
	"user/pkg/config"
	"user/pkg/zjwt"
)

func (c *LoginController) getAcrossSubdomainCookie(accessToken string) []*http.Cookie {
	cnf := config.GetConfig()
	rootDomain := cnf.RootDomain
	accessTokenCookie := &http.Cookie{
		Name:    "sso_0voice_access_token",
		Value:   accessToken,
		Path:    "/",
		Domain:  rootDomain,
		Expires: time.Now().Add(time.Second * time.Duration(zjwt.EXPIRES_IN)),
	}
	localAccessTokenCookie := &http.Cookie{
		Name:    "sso_0voice_access_token",
		Value:   accessToken,
		Path:    "/",
		Expires: time.Now().Add(time.Second * time.Duration(zjwt.EXPIRES_IN)),
	}
	return []*http.Cookie{accessTokenCookie, localAccessTokenCookie}
}
