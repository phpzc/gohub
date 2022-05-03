package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

//处理请求数据和表单验证

type SignupPhoneExistRequest struct {
	Phone string `json:"phone,omitempty" valid:"phone"`
}

func ValidateSignupPhoneExist(data interface{}, c *gin.Context) map[string][]string {
	//自定义验证规则
	rules := govalidator.MapData{
		"phone": []string{"required", "digits:11"},
	}

	//自定义验证出错提示
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项,参数名称phone",
			"digits:手机号长度必须为11位的数字",
		},
	}

	//配置初始化
	opts := govalidator.Options{

		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	//开始验证
	return govalidator.New(opts).ValidateStruct()
}

type SignupEmailExistRequest struct {
	Email string `json:"email,omitempty" valid:"email"`
}

//验证邮箱是否已注册
func ValidateSignupEmailExist(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		"email": []string{
			"required",
			"min:4",
			"max:30",
			"email",
		},
	}

	messages := govalidator.MapData{
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需大于4",
			"max:Email 长度需小于30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
		},
	}

	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	return govalidator.New(opts).ValidateStruct()
}
