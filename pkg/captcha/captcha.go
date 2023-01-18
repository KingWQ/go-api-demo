package captcha

import (
	"github.com/mojocn/base64Captcha"
	"go-api-demo/pkg/app"
	"go-api-demo/pkg/config"
	"go-api-demo/pkg/redis"
	"sync"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

// once 确保 internalCaptcha 对象只初始化一次
var once sync.Once

// internalCaptcha 内部使用的 Captcha对象
var internalCaptcha *Captcha

// NewCaptcha 单例模式获取
func NewCaptcha() *Captcha {
	once.Do(func() {
		//1，初始化Captcha对象
		internalCaptcha = &Captcha{}

		//2，使用全局Redis对象，并配置存储Key前缀
		store := RedisStore{
			RedisClient: redis.Redis,
			KeyPrefix:   config.GetString("app.name") + ":captcha",
		}

		//3，配置base64Captcha驱动信息
		driver := base64Captcha.NewDriverDigit(
			config.GetInt("captcha.height"),      // 宽
			config.GetInt("captcha.width"),       // 高
			config.GetInt("captcha.length"),      // 长度
			config.GetFloat64("captcha.maxskew"), // 数字的最大倾斜角度
			config.GetInt("captcha.dotcount"),    // 图片背景里的混淆点数量
		)

		//4，实例化base64Captcha 并赋值给内部使用的 internalCaptcha对象
		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, &store)

	})
	return internalCaptcha
}

// GenerateCaptcha 生成图片验证码
func (c *Captcha) GenerateCaptcha() (id string, b64s string, err error) {
	return c.Base64Captcha.Generate()
}

// VerifyCaptcha 验证验证码是否正确
func (c *Captcha) VerifyCaptcha(id, answer string) (match bool) {
	// 方便本地和 API 自动测试
	if !app.IsProduction() && id == config.GetString("captcha.testing_key") {
		return true
	}
	
	// 第三个参数是验证后是否删除，我们选择 false
	// 这样方便用户多次提交，防止表单提交错误需要多次输入图片验证码
	return c.Base64Captcha.Verify(id, answer, false)
}
