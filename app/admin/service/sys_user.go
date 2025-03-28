package service

import (
	"errors"
	"quanta-admin/app/admin/models"
	"quanta-admin/app/admin/service/dto"
	"quanta-admin/common/utils"

	"github.com/pquerna/otp/totp"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type SysUser struct {
	service.Service
}

// GetPage 获取SysUser列表
func (e *SysUser) GetPage(c *dto.SysUserGetPageReq, p *actions.DataPermission, list *[]models.SysUser, count *int64) error {
	var err error
	var data models.SysUser

	err = e.Orm.Debug().
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// Get 获取SysUser对象
func (e *SysUser) Get(d *dto.SysUserById, p *actions.DataPermission, model *models.SysUser) error {
	var data models.SysUser

	err := e.Orm.Model(&data).Debug().
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// Get2FACode 为sysuser生成2fa绑定的二维码
func (e *SysUser) Get2FACode(userId *int, resp *dto.GetTwoFaCodeResp) error {
	data := models.SysUser{}

	err := e.Orm.Model(&data).Debug().
		First(&data, userId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if data.ActiveTwoFa {
		err := errors.New("当前用户已启用2fa")
		e.Log.Errorf("当前用户已经启用2fa: %s", err)
		return err
	}

	faSecret, faCode, err := utils.Generate2FA(data.Username)
	if err != nil {
		e.Log.Errorf("生成2fa密钥失败: %s", err)
		return err
	}

	hashedSecret, err := utils.Encrypt(faSecret)
	if err != nil {
		e.Log.Errorf("2fa密钥加密失败: %s", err)
		return err
	}
	err = e.Orm.Model(&data).Update("two_fa_secret", hashedSecret).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	resp.TwoFaSecret = faSecret
	resp.TwoFaCodeUrl = faCode

	return nil
}

// Verify2FACode 为sysuser绑定2fa验证code
func (e *SysUser) Verify2FACode(userId *int, req *dto.BindTwoFaVerifyRequest) error {
	data := models.SysUser{}
	e.Log.Infof("userId:%d, req:%+v", userId, req)
	err := e.Orm.Model(&data).Debug().
		First(&data, userId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if data.ActiveTwoFa {
		err := errors.New("当前用户已启用2fa，不能重复绑定")
		e.Log.Errorf("当前用户已经启用2fa: %s", err)
		return err
	}

	// 从数据库中获取哈希后的密钥
	encryptSecret := data.TwoFaSecret
	decryptSecret, err := utils.Decrypt(encryptSecret)
	if err != nil {
		e.Log.Errorf("2fa密钥解密失败: %s", err)
		return err
	}
	valid := totp.Validate(req.TwoFaCode, decryptSecret)
	if !valid {
		e.Log.Errorf("2fa验证失败")
		err = errors.New("2fa验证失败")
		return err
	}
	e.Log.Info("2fa验证结果", valid)

	err = e.Orm.Model(&data).Debug().
		Where("user_id = ?", userId).
		Update("active_two_fa", true).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	return nil
}

// Unbind2FACode 为SysUser解绑2fa验证code
func (e *SysUser) Unbind2FACode(userId *int, req *dto.UnBindTwoFaVerifyRequest) error {
	data := models.SysUser{}
	err := e.Orm.Model(&data).Debug().
		First(&data, userId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if !data.ActiveTwoFa {
		err := errors.New("当前用户未绑定2fa")
		e.Log.Errorf("当前用户未绑定2fa: %s", err)
		return err
	}

	// 从数据库中获取哈希后的密钥
	encryptSecret := data.TwoFaSecret
	decryptSecret, err := utils.Decrypt(encryptSecret)
	if err != nil {
		e.Log.Errorf("2fa密钥解密失败: %s", err)
		return err
	}
	valid := totp.Validate(req.TwoFaCode, decryptSecret)
	if !valid {
		e.Log.Errorf("2fa验证失败")
		err = errors.New("2fa验证失败")
		return err
	}
	e.Log.Info("2fa验证结果", valid)

	err = e.Orm.Model(&data).Debug().
		Where("user_id = ?", userId).
		Update("active_two_fa", false).
		Update("two_fa_secret", "").Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	return nil
}

// Insert 创建SysUser对象
func (e *SysUser) Insert(c *dto.SysUserInsertReq) error {
	var err error
	var data models.SysUser
	var i int64
	err = e.Orm.Model(&data).Where("username = ?", c.Username).Count(&i).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if i > 0 {
		err := errors.New("用户名已存在！")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// Update 修改SysUser对象
func (e *SysUser) Update(c *dto.SysUserUpdateReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	c.Generate(&model)
	update := e.Orm.Model(&model).Where("user_id = ?", &model.UserId).Omit("password", "salt").Updates(&model)
	if err = update.Error; err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if update.RowsAffected == 0 {
		err = errors.New("update userinfo error")
		log.Warnf("db update error")
		return err
	}
	return nil
}

// UpdateAvatar 更新用户头像
func (e *SysUser) UpdateAvatar(c *dto.UpdateSysUserAvatarReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	err = e.Orm.Table(model.TableName()).Where("user_id =? ", c.UserId).Updates(c).Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	return nil
}

// UpdateStatus 更新用户状态
func (e *SysUser) UpdateStatus(c *dto.UpdateSysUserStatusReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	err = e.Orm.Table(model.TableName()).Where("user_id =? ", c.UserId).Updates(c).Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	return nil
}

// ResetPwd 重置用户密码
func (e *SysUser) ResetPwd(c *dto.ResetSysUserPwdReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("At Service ResetSysUserPwd error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	c.Generate(&model)
	err = e.Orm.Omit("username", "nick_name", "phone", "role_id", "avatar", "sex").Save(&model).Error
	if err != nil {
		e.Log.Errorf("At Service ResetSysUserPwd error: %s", err)
		return err
	}
	return nil
}

// Remove 删除SysUser
func (e *SysUser) Remove(c *dto.SysUserById, p *actions.DataPermission) error {
	var err error
	var data models.SysUser

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Error found in  RemoveSysUser : %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// UpdatePwd 修改SysUser对象密码
func (e *SysUser) UpdatePwd(id int, oldPassword, newPassword string, p *actions.DataPermission) error {
	var err error

	if newPassword == "" {
		return nil
	}
	c := &models.SysUser{}

	err = e.Orm.Model(c).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Select("UserId", "Password", "Salt").
		First(c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("无权更新该数据")
		}
		e.Log.Errorf("db error: %s", err)
		return err
	}
	var ok bool
	ok, err = pkg.CompareHashAndPassword(c.Password, oldPassword)
	if err != nil {
		e.Log.Errorf("CompareHashAndPassword error, %s", err.Error())
		return err
	}
	if !ok {
		err = errors.New("incorrect Password")
		e.Log.Warnf("user[%d] %s", id, err.Error())
		return err
	}
	c.Password = newPassword
	db := e.Orm.Model(c).Where("user_id = ?", id).
		Select("Password", "Salt").
		Updates(c)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("set password error")
		log.Warnf("db update error")
		return err
	}
	return nil
}

func (e *SysUser) GetProfile(c *dto.SysUserById, user *models.SysUser, roles *[]models.SysRole) error {
	err := e.Orm.First(user, c.GetId()).Error
	if err != nil {
		return err
	}

	// 解析 RoleIds 字符串为 []int
	roleIds, err := user.ParseRoleIds()
	if err != nil {
		return err
	}

	// 查询所有角色信息
	if len(roleIds) > 0 {
		err = e.Orm.Where("role_id IN (?)", roleIds).Find(roles).Error
		if err != nil {
			return err
		}
	} else {
		// 如果没有角色ID，则返回空数组
		*roles = []models.SysRole{}
	}

	return nil
}
