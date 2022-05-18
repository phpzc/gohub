package verifycode

import (
	"fmt"
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
	"gohub/pkg/mail"
	"gohub/pkg/redis"
	"gohub/pkg/sms"
	"strings"
	"sync"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once
var internalVerifyCode *VerifyCode

// NewVerifyCode 单例模式获取
func NewVerifyCode() *VerifyCode {
	once.Do(func() {

		internalVerifyCode = &VerifyCode{
			Store: &RedisStore{
				RedisClient: redis.Redis,
				//增加前缀保存数据库整洁
				KeyPrefix: config.GetString("app.name") + ":verifycode:",
			},
		}
	})

	return internalVerifyCode
}

// 生成验证码 放入redis中
func (vc *VerifyCode) generateVerifyCode(key string) string {

	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))

	if app.IsLocal() {
		code = config.GetString("verifycode.debug_code")
	}

	vc.Store.Set(key, code)

	return code

}

// SendSMS 发送短信验证码，调用示例：
//         verifycode.NewVerifyCode().SendSMS(request.Phone)
func (vc *VerifyCode) SendSMS(phone string) bool {

	// 生成验证码
	code := vc.generateVerifyCode(phone)

	// 方便本地和 API 自动测试
	if !app.IsProduction() && strings.HasPrefix(phone, config.GetString("verifycode.debug_phone_prefix")) {
		return true
	}

	// 发送短信
	return sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.aliyun.template_code"),
		Data:     map[string]string{"code": code},
	})
}

// CheckAnswer 检查用户提交的验证码是否正确，key 可以是手机号或者 Email
func (vc *VerifyCode) CheckAnswer(key string, answer string) bool {

	logger.DebugJSON("验证码", "检查验证码", map[string]string{key: answer})

	// 方便开发，在非生产环境下，具备特殊前缀的手机号和 Email后缀，会直接验证成功
	if !app.IsProduction() &&
		(strings.HasSuffix(key, config.GetString("verifycode.debug_email_suffix")) ||
			strings.HasPrefix(key, config.GetString("verifycode.debug_phone_prefix"))) {
		return true
	}

	return vc.Store.Verify(key, answer, false)
}

//SendEmail 发送邮件验证码
func (vc *VerifyCode) SendEmail(email string) error {
	code := vc.generateVerifyCode(email)

	if !app.IsProduction() && strings.HasSuffix(email, config.GetString("verifycode.debug_email_suffix")) {
		return nil
	}

	content := fmt.Sprintf("<h1>您的 Email 验证码是 %v </h1>", code)

	mail.NewMailer().Send(mail.Email{
		From: mail.From{
			Address: config.GetString("mail.from.address"),
			Name:    config.GetString("mail.from.name"),
		},

		To:      []string{email},
		Subject: "Email 验证码",
		HTML:    []byte(content),
	})

	return nil
}
