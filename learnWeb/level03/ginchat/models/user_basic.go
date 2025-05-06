package models

import (
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	ID            uint           `json:"id" example:"1"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Name          string
	Password      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	HeartBeatTime uint64
	LoginOutTime  uint64
	IsLogout      bool
	DeviceInfo    string
	Salt          string `gorm:"column:Salt"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	db := utils.GetDB()
	db.Find(&data)
	return data
}

func CreateUser(user UserBasic) *gorm.DB {
	db := utils.GetDB()
	db.Create(&user)
	return db
}

func DeleteUser(id int) *gorm.DB {

	db := utils.GetDB()
	//Unscoped()是硬删除，去掉Unscoped()是软删除，因为model里面添加了依赖提供的字段
	db.Unscoped().Where("id = ?", id).Delete(&UserBasic{})
	return db
}

func UpdateUser(user *UserBasic) *gorm.DB {
	db := utils.GetDB()
	db.Exec("update user_basic set updated_at = ?,phone = ?,name =?,password =? where id = ?", time.Now(), user.Phone, user.Name, user.Password, user.ID)
	return db
}

func GetUserByName(name string) *UserBasic {
	db := utils.GetDB()
	var user UserBasic
	db.Where("name = ?", name).First(&user)
	return &user
}
