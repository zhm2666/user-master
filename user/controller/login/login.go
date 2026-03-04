package login

import (
	"github.com/gin-gonic/gin"
	redis2 "github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"net/http"
	"user/data"
	"user/pkg/config"
	"user/pkg/constants"
	gitlab2 "user/pkg/gitlab"
	"user/pkg/log"
	"user/pkg/storage"
	"user/pkg/utils"
	"user/pkg/zerror"
)

type LoginController struct {
	log         log.ILogger
	config      *config.Config
	sf          storage.StorageFactory
	data        data.IData
	redisClient *redis2.Client
}

func NewLoginController(log log.ILogger, config *config.Config, redisClient *redis2.Client, sf storage.StorageFactory, data data.IData) *LoginController {
	return &LoginController{
		log:         log,
		config:      config,
		redisClient: redisClient,
		data:        data,
		sf:          sf,
	}
}

func (c *LoginController) GetLoginMethods(ctx *gin.Context) {
	sys := ctx.DefaultQuery("sys", string(constants.SYS_MEDIAHUB))
	if sys == "" {
		err := zerror.NewByMsg("请指定需要登录的系统")
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	gitlab := gitlab2.NewGitlabOAuth(c.config.Gitlab.Domain, c.config.Gitlab.ClientID, c.config.Gitlab.ClientSecret, c.log)
	gitlabConf := gitlab.GetOauth2Config(c.config.Gitlab.RedirectUri, map[string]string{"sys": sys}, []string{"read_user"})
	gitlabAuthUrl := gitlabConf.AuthCodeURL(utils.GenerateRandomString(8), oauth2.SetAuthURLParam("grant_type", "authorization_code"))

	qrcode, err := c.getWxMpQrCode()
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"gitlab":    gitlabAuthUrl,
		"wx_qrcode": qrcode,
	})
	return
}
