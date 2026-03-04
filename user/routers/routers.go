package routers

import (
	"github.com/gin-gonic/gin"
	"user/controller/login"
)

func InitLoginRouters(api *gin.RouterGroup, c *login.LoginController) {
	v1 := api.Group("/v1")
	group := v1.Group("/login")
	group.GET("/official/receive", c.WxCheckSignature)
	group.POST("/official/receive", c.WxReceiveMessage)
	group.GET("/methods", c.GetLoginMethods)
	group.GET("/official/callback", c.WxCallBackByTicket)
	group.GET("/gitlab/redirect", c.OAuthCallback)
	group.GET("/check/auth", c.CheckLogin)
}
