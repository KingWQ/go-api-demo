// Package user 存放用户 Model 相关逻辑
package user

import (
	"go-api-demo/app/models"
	"go-api-demo/pkg/database"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"column:name;type:varchar(60)" json:"name,omitempty"`
	Email    string `gorm:"column:email;type:varchar(60)" json:"-"`
	Phone    string `gorm:"column:phone;type:varchar(20)" json:"-"`
	Password string `gorm:"column:password;type:char(32)" json:"-"`

	models.CommonTimestampsField
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (userModel *User) Create() {
	database.DB.Create(&userModel)
}
