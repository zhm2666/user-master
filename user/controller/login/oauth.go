package login

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"net/http"
	"time"
	"user/data"
	gitlab2 "user/pkg/gitlab"
	"user/pkg/zerror"
	"user/pkg/zjwt"
)

func (c *LoginController) OAuthCallback(ctx *gin.Context) {
	code := ctx.DefaultQuery("code", "")
	if code == "" {
		err := zerror.NewByMsg("授权失败")
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	sys := ctx.DefaultQuery("sys", "")
	if sys == "" {
		err := zerror.NewByMsg("请指定业务系统")
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	gitlab := gitlab2.NewGitlabOAuth(c.config.Gitlab.Domain, c.config.Gitlab.ClientID, c.config.Gitlab.ClientSecret, c.log)
	gitlabConf := gitlab.GetOauth2Config(c.config.Gitlab.RedirectUri, map[string]string{"sys": sys}, []string{"read_user"})
	token, err := gitlabConf.Exchange(context.Background(), code, oauth2.SetAuthURLParam("grant_type", "authorization_code"))
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	gitlabUser, err := gitlab.GetUser(token)
	gitlabUserData := c.data.NewGitlabUserData()
	gitlabEntity, err := gitlabUserData.GetByGitlabID(gitlabUser.ID)
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if gitlabEntity == nil {
		gitlabEntity = &data.GitlabUser{
			GitlabID:  gitlabUser.ID,
			UserName:  gitlabUser.UserName,
			Name:      gitlabUser.Name,
			Email:     gitlabUser.Email,
			AvatarUrl: gitlabUser.AvatarUrl,
			CreateAt:  time.Now().Unix(),
		}
		err = gitlabUserData.AddUser(gitlabEntity)
		if err != nil {
			c.log.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
	}
	userEntity, err := c.data.NewUserData().GetByID(gitlabEntity.UserID)
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	hs := zjwt.NewHs(c.config.Jwt.HashKey, zjwt.HS256)
	now := time.Now()
	userClaims := zjwt.UserClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(zjwt.EXPIRES_IN))),
		},
		UserID: userEntity.ID,
		Name:   userEntity.Name,
		Avatar: userEntity.AvatarUrl,
	}
	accessToken, err := hs.Sign(userClaims)
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	cookies := c.getAcrossSubdomainCookie(accessToken)
	for _, cookie := range cookies {
		http.SetCookie(ctx.Writer, cookie)
	}
	redirectUrl := c.config.InternalSystemEntry[sys]
	ctx.Redirect(http.StatusFound, redirectUrl)
}
