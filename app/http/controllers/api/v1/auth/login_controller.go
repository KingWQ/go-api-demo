package auth

import (
	"github.com/gin-gonic/gin"
	v1 "go-api-demo/app/http/controllers/api/v1"
	"go-api-demo/app/requests"
	"go-api-demo/pkg/auth"
	"go-api-demo/pkg/jwt"
	"go-api-demo/pkg/response"
)

// LoginController 用户控制器
type LoginController struct {
	v1.BaseAPIController
}

// LoginByPhone 手机登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
		return
	}

	user, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		response.Error(c, err, "账号不存在或密码错误")
	} else {
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
		response.JSON(c, gin.H{
			"token": token,
		})
	}

}
