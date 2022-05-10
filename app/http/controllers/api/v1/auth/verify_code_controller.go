package auth

import (
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/pkg/captcha"
	"gohub/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VerifyCodeController struct {
	v1.BaseAPIController
}

//显示图片验证码
func (vc *VerifyCodeController) ShowCaptcha(c *gin.Context) {

	//生成验证码
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	//记录错误日志
	logger.LogIf(err)

	//返回给用户
	c.JSON(http.StatusOK, gin.H{
		"captch_id":    id,
		"cpatch_image": b64s,
	})
}
