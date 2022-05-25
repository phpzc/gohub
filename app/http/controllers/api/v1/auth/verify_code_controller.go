package auth

import (
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/requests"
	"gohub/pkg/captcha"
	"gohub/pkg/logger"
	"gohub/pkg/response"
	"gohub/pkg/verifycode"

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
	response.JSON(c, gin.H{
		"captch_id":    id,
		"captch_image": b64s,
	})
}

// 发送Email验证码
func (vc *VerifyCodeController) SendUsingEmail(c *gin.Context) {

	//验证表单
	request := requests.VerifyCodeEmailRequest{}

	if ok := requests.Validate(c, &request, requests.VerifyCodeEmail); !ok {
		return
	}

	//发送
	err := verifycode.NewVerifyCode().SendEmail(request.Email)

	if err != nil {
		response.Abort500(c, "发送 Email 验证码失败")
	} else {
		response.Success(c)
	}
}

// SendUsingPhone 发送手机验证码
func (vc *VerifyCodeController) SendUsingPhone(c *gin.Context) {

	// 1. 验证表单
	request := requests.VerifyCodePhoneRequest{}
	if ok := requests.Validate(c, &request, requests.VerifyCodePhone); !ok {
		return
	}

	// 2. 发送 SMS
	if ok := verifycode.NewVerifyCode().SendSMS(request.Phone); !ok {
		response.Abort500(c, "发送短信失败~")
	} else {
		response.Success(c)
	}
}
