package verifycode

import (
	"fmt"
	"go-api-demo/pkg/app"
	"go-api-demo/pkg/config"
	"go-api-demo/pkg/helpers"
	"go-api-demo/pkg/logger"
	"go-api-demo/pkg/mail"
	"go-api-demo/pkg/redis"
	"go-api-demo/pkg/sms"
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
				KeyPrefix:   config.GetString("app.name") + ":verifycode:",
			},
		}
	})

	return internalVerifyCode
}

// SendSMS 发送短信验证码，调用示例：verifycode.NewVerifyCode().SendSMS(request.Phone)
func (vc *VerifyCode) SendSMS(phone string) bool {
	//1. 生成短信验证码
	code := vc.generateVerifyCode(phone)

	//2. 方便本地和api自动测试
	if !app.IsProduction() && strings.HasPrefix(phone, config.GetString("verifycode.debug_phone_prefix")) {
		return true
	}

	//3. 发送短信
	return sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.aliyun.template_code"),
		Data:     map[string]string{"code": code},
	})
}

// SendEmail 发送邮件验证码，调用示例：verifycode.NewVerifyCode().SendEmail(request.Phone)
func (vc *VerifyCode) SendEmail(email string) error {
	//1. 生成短信验证码
	code := vc.generateVerifyCode(email)

	//2. 方便本地和api自动测试
	if !app.IsProduction() && strings.HasSuffix(email, config.GetString("verifycode.debug_email_prefix")) {
		return nil
	}

	content := fmt.Sprintf("<h1>您的 Email 验证码是 %v </h1>", code)
	//3. 发送邮件
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

// CheckAnswer 检查用户提交的验证码是否正确，key可以是手机号或email
func (vc *VerifyCode) CheckAnswer(key, answer string) bool {
	logger.DebugJSON("验证码", "检查验证码", map[string]string{key: answer})

	// 方便开发，在非生产环境下，具备特殊前缀的手机号和 Email后缀，会直接验证成功
	if !app.IsProduction() &&
		(strings.HasSuffix(key, config.GetString("verifycode.debug_phone_prefix")) ||
			strings.HasSuffix(key, config.GetString("verifycode.debug_email_suffix"))) {
		return true
	}

	return vc.Store.Verify(key, answer, false)
}

// generateVerifyCode 生成验证码，并放置于Redis中
func (vc *VerifyCode) generateVerifyCode(key string) string {
	//1. 生成随机验证码
	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))

	//2. 为方便开发，本地环境使用固定验证码
	if app.IsLocal() {
		code = config.GetString("verifycode.debug_code")
	}

	logger.DebugJSON("验证码", "生成验证码", map[string]string{key: code})

	//3. 将验证码和key(邮箱或手机号) 存放到Redis中并设置过期时间
	vc.Store.Set(key, code)
	return code
}
