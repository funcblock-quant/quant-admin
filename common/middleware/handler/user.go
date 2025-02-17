package handler

import (
	"gorm.io/gorm"
	"quanta-admin/common/models"
)

type SysUser struct {
	UserId             int    `gorm:"primaryKey;autoIncrement;comment:编码"  json:"userId"`
	Username           string `json:"username" gorm:"size:64;comment:用户名"`
	Password           string `json:"-" gorm:"size:128;comment:密码"`
	NickName           string `json:"nickName" gorm:"size:128;comment:昵称"`
	Phone              string `json:"phone" gorm:"size:11;comment:手机号"`
	RoleId             int    `json:"roleId" gorm:"size:20;comment:角色ID"`
	Salt               string `json:"-" gorm:"size:255;comment:加盐"`
	Avatar             string `json:"avatar" gorm:"size:255;comment:头像"`
	Sex                string `json:"sex" gorm:"size:255;comment:性别"`
	ActiveTwoFa        bool   `json:"-" gorm:"comment:是否开启"`
	TwoFaSecret        string `json:"-" gorm:"size:255;comment:2fa密钥"`
	TwoFaRecoverSecret string `json:"-" gorm:"size:255;comment:2fa一次性恢复密钥"`
	Email              string `json:"email" gorm:"size:128;comment:邮箱"`
	Remark             string `json:"remark" gorm:"size:255;comment:备注"`
	Status             string `json:"status" gorm:"size:4;comment:状态"`
	RoleIds            []int  `json:"roleIds" gorm:"-"`
	//Dept     *SysDept `json:"dept"`
	models.ControlBy
	models.ModelTime
}

func (*SysUser) TableName() string {
	return "sys_user"
}

func (e *SysUser) AfterFind(_ *gorm.DB) error {
	e.RoleIds = []int{e.RoleId}
	return nil
}
