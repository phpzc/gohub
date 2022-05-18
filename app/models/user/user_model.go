package user

import (
	"gohub/app/models"
	"gohub/pkg/database"
)

type User struct {
	models.BaseModel

	//因为我们不希望将敏感信息输出给用户，所以这里 Email 、Phone 、Password 后面设置了 json:"-" ，这是在指示 JSON 解析器忽略字段
	Name     string `json:"name,omitempty"`
	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`

	models.CommonTimeStampsField
}

func (userModel *User) Create() {
	database.DB.Create(&userModel)
}
