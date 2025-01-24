package handler

import (
	"errors"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"quanta-admin/common/utils"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

func (u *Login) GetUser(tx *gorm.DB) (user SysUser, role SysRole, err error) {
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
	err = tx.Table("sys_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	if err != nil {
		log.Errorf("get role error, %s", err.Error())
		return
	}
	return
}

type TwoFALogin struct {
	Username string `form:"Username" json:"username" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
}

func (u *TwoFALogin) GetUser(tx *gorm.DB) (user SysUser, role SysRole, err error) {
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
	err = tx.Table("sys_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	log.Info("role", role)
	if err != nil {
		log.Errorf("get role error, %s", err.Error())
		return
	}
	return
}
