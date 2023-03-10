// Package user 存放用户 Model 相关逻辑
package user

import (
	"go-api-demo/app/models"
	"go-api-demo/pkg/database"
	"go-api-demo/pkg/hash"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"column:name;type:varchar(60)" json:"name,omitempty"`
	Email    string `gorm:"column:email;type:varchar(60)" json:"-"`
	Phone    string `gorm:"column:phone;type:varchar(20)" json:"-"`
	Password string `gorm:"column:password;type:char(60)" json:"-"`

	models.CommonTimestampsField
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (userModel *User) Create() {
	database.DB.Create(&userModel)
}

// ComparePassword 密码是否正确
func (userModel *User) ComparePassword(_password string) bool {
	return hash.BcryptCheck(_password, userModel.Password)
}

func (userModel *User) Save() (rowsAffected int64) {
	result := database.DB.Save(&userModel)
	return result.RowsAffected
}
