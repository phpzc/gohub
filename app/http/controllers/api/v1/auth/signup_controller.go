package auth

import (
	"fmt"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupController struct {
	v1.BaseAPIController
}

//检测手机号是否被注册
func (sc *SignupController) IsPhoneExist(c *gin.Context) {

	//请求对象
	// type PhoneExistRequest struct {
	// 	Phone string `json:"phone"`
	// }

	// request := PhoneExistRequest{}

	//初始化请求对象
	request := requests.SignupPhoneExistRequest{}

	//解析JSON请求
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})

		fmt.Println(err.Error())
		return
	}

	//表单验证
	errs := requests.ValidateSignupPhoneExist(&request, c)

	if len(errs) > 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"erros": errs,
		})
		return
	}

	//检查数据库 并返回响应
	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}
