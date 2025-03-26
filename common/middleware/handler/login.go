package handler

import (
	"errors"
	"quanta-admin/common/utils"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

func (u *Login) GetUser(tx *gorm.DB) (user SysUser, roles []SysRole, err error) {
	err = tx.Table("sys_user").Where("username = ?  and status = '2'", u.Username).First(&user).Error
	if err != nil {
		log.Errorf("get user error, %s", err.Error())
		return
	}
	_, err = pkg.CompareHashAndPassword(user.Password, u.Password)
	if err != nil {
		log.Errorf("user login error, %s", err.Error())
		return
	}

	// 解析 RoleIds 字符串为 []int
	roleIds, err := user.ParseRoleIds()
	if err != nil {
		log.Errorf("parse role ids error, %s", err.Error())
		return
	}

	// 查询所有角色信息
	if len(roleIds) > 0 {
		err = tx.Table("sys_role").Where("role_id IN (?)", roleIds).Find(&roles).Error
		if err != nil {
			log.Errorf("get roles error, %s", err.Error())
			return
		}
	} else {
		// 如果没有角色 ID，则返回空数组
		roles = []SysRole{}
	}
	return
}

type TwoFALogin struct {
	Username string `form:"Username" json:"username" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
}

func (u *TwoFALogin) GetUser(tx *gorm.DB) (user SysUser, roles []SysRole, err error) {
	err = tx.Table("sys_user").Where("username = ?  and status = '2'", u.Username).First(&user).Error
	if err != nil {
		log.Errorf("get user error, %s", err.Error())
		return
	}
	//2fa 验证
	encryptSecret := user.TwoFaSecret
	decryptSecret, err := utils.Decrypt(encryptSecret)
	if err != nil {
		log.Errorf("2fa密钥解密失败: %s", err)
		return
	}
	valid := totp.Validate(u.Code, decryptSecret)
	if !valid {
		log.Errorf("2fa验证失败")
		err = errors.New("2fa验证失败")
	}
	log.Info("2fa验证结果", valid)

	if err != nil {
		log.Errorf("user login error, %s", err.Error())
		return
	}
	// 解析 RoleIds 字符串为 []int
	roleIds, err := user.ParseRoleIds()
	if err != nil {
		log.Errorf("parse role ids error, %s", err.Error())
		return
	}

	// 查询所有角色信息
	if len(roleIds) > 0 {
		err = tx.Table("sys_role").Where("role_id IN (?)", roleIds).Find(&roles).Error
		if err != nil {
			log.Errorf("get roles error, %s", err.Error())
			return
		}
	} else {
		// 如果没有角色 ID，则返回空数组
		roles = []SysRole{}
	}
	return
}
