package user

import "go-api-demo/pkg/database"

// IsEmailExist 判断Email是否注册
func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// IsPhoneExist 判断手机号是否注册
func IsPhoneExist(phone string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

func GetByPhone(phone string) (userModel User) {
	database.DB.Where("phone = ?", phone).First(&userModel)
	return
}

func GetByMulti(loginID string) (userModel User) {
	database.DB.Where("phone = ?", loginID).
		Or("email = ?", loginID).
		Or("name = ?", loginID).
		First(&userModel)
	return
}

func Get(idStr string) (userModel User) {
	database.DB.Where("id", idStr).First(&userModel)
	return
}
