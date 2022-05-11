package auth

import (
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

type SignupController struct {
	v1.BaseAPIController
}

//检测手机号是否被注册
func (sc *SignupController) IsPhoneExist(c *gin.Context) {

	//初始化请求对象
	request := requests.SignupPhoneExistRequest{}

	//请求参数 与 表单验证
	if ok := requests.Validate(c, &request, requests.ValidateSignupPhoneExist); !ok {
		return
	}

	//检查数据库 并返回响应
	response.JSON(c, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}

//检测邮箱是否已被注册
func (sc *SignupController) IsEmailExist(c *gin.Context) {
	//初始化请求对象
	request := requests.SignupEmailExistRequest{}

	//请求参数 与 表单验证
	if ok := requests.Validate(c, &request, requests.ValidateSignupEmailExist); !ok {
		return
	}

	//检查数据库 并返回响应
	response.JSON(c, gin.H{
		"exist": user.IsEmailExist(request.Email),
	})
}
