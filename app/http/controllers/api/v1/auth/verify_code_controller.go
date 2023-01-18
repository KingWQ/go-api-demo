package auth

import (
	"github.com/gin-gonic/gin"
	v1 "go-api-demo/app/http/controllers/api/v1"
	"go-api-demo/pkg/captcha"
	"go-api-demo/pkg/logger"
	"net/http"
)

// VerifyCodeController 用户控制器
type VerifyCodeController struct {
	v1.BaseAPIController
}

// ShowCaptcha 显示图片验证码
func (vc *VerifyCodeController) ShowCaptcha(c *gin.Context) {
	//1.生成验证码
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()

	//2. 记录错误日志
	logger.LogIf(err)

	//3. 返回给用户
	c.JSON(http.StatusOK, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}