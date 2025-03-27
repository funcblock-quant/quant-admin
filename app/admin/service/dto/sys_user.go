package dto

import (
	"quanta-admin/app/admin/models"
	"strconv"

	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type SysUserGetPageReq struct {
	dto.Pagination `search:"-"`
	UserId         int    `form:"userId" search:"type:exact;column:user_id;table:sys_user" comment:"用户ID"`
	Username       string `form:"username" search:"type:contains;column:username;table:sys_user" comment:"用户名"`
	NickName       string `form:"nickName" search:"type:contains;column:nick_name;table:sys_user" comment:"昵称"`
	Phone          string `form:"phone" search:"type:contains;column:phone;table:sys_user" comment:"手机号"`
	RoleId         string `form:"roleId" search:"type:exact;column:role_id;table:sys_user" comment:"角色ID"`
	Sex            string `form:"sex" search:"type:exact;column:sex;table:sys_user" comment:"性别"`
	Email          string `form:"email" search:"type:contains;column:email;table:sys_user" comment:"邮箱"`
	PostId         string `form:"postId" search:"type:exact;column:post_id;table:sys_user" comment:"岗位"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_user" comment:"状态"`
	SysUserOrder
}

type SysUserOrder struct {
	UserIdOrder    string `search:"type:order;column:user_id;table:sys_user" form:"userIdOrder"`
	UsernameOrder  string `search:"type:order;column:username;table:sys_user" form:"usernameOrder"`
	StatusOrder    string `search:"type:order;column:status;table:sys_user" form:"statusOrder"`
	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_user" form:"createdAtOrder"`
}

func (m *SysUserGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ResetSysUserPwdReq struct {
	UserId   int    `json:"userId" comment:"用户ID" vd:"$>0"` // 用户ID
	Password string `json:"password" comment:"密码" vd:"len($)>0"`
	common.ControlBy
}

func (s *ResetSysUserPwdReq) GetId() interface{} {
	return s.UserId
}

func (s *ResetSysUserPwdReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Password = s.Password
}

type UpdateSysUserAvatarReq struct {
	UserId int    `json:"userId" comment:"用户ID" vd:"len($)>0"` // 用户ID
	Avatar string `json:"avatar" comment:"头像" vd:"len($)>0"`
	common.ControlBy
}

func (s *UpdateSysUserAvatarReq) GetId() interface{} {
	return s.UserId
}

func (s *UpdateSysUserAvatarReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Avatar = s.Avatar
}

type UpdateSysUserStatusReq struct {
	UserId int    `json:"userId" comment:"用户ID" vd:"$>0"` // 用户ID
	Status string `json:"status" comment:"状态" vd:"len($)>0"`
	common.ControlBy
}

func (s *UpdateSysUserStatusReq) GetId() interface{} {
	return s.UserId
}

func (s *UpdateSysUserStatusReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Status = s.Status
}

type SysUserInsertReq struct {
	UserId   int    `json:"userId" comment:"用户ID"` // 用户ID
	Username string `json:"username" comment:"用户名" vd:"len($)>0"`
	Password string `json:"password" comment:"密码"`
	NickName string `json:"nickName" comment:"昵称" vd:"len($)>0"`
	Phone    string `json:"phone" comment:"手机号" vd:"len($)>0"`
	RoleIds  []int  `json:"roleIds" comment:"角色ID"`
	Avatar   string `json:"avatar" comment:"头像"`
	Sex      string `json:"sex" comment:"性别"`
	Email    string `json:"email" comment:"邮箱" vd:"len($)>0,email"`
	Remark   string `json:"remark" comment:"备注"`
	Status   string `json:"status" comment:"状态" vd:"len($)>0" default:"1"`
	common.ControlBy
}

func (s *SysUserInsertReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Username = s.Username
	model.Password = s.Password
	model.NickName = s.NickName
	model.Phone = s.Phone
	model.Avatar = s.Avatar
	model.Sex = s.Sex
	model.Email = s.Email
	model.Remark = s.Remark
	model.Status = s.Status
	model.CreateBy = s.CreateBy

	// 拼接 RoleIds 字符串
	roleIdsStr := ""
	for i, roleId := range s.RoleIds {
		roleIdsStr += strconv.Itoa(roleId)
		if i < len(s.RoleIds)-1 {
			roleIdsStr += ","
		}
	}
	model.RoleIds = roleIdsStr
}

func (s *SysUserInsertReq) GetId() interface{} {
	return s.UserId
}

type GetTwoFaCodeResp struct {
	TwoFaSecret  string `json:"twoFaSecret" comment:"2fa密钥"`
	TwoFaCodeUrl string `json:"twoFaCodeUrl" comment:"2fa二维码"`
}

type BindTwoFaVerifyRequest struct {
	TwoFaCode string `json:"twoFaCode" comment:"2fa验证码"`
}

type UnBindTwoFaVerifyRequest struct {
	TwoFaCode string `json:"twoFaCode" comment:"2fa验证码"`
}

type SysUserUpdateReq struct {
	UserId   int    `json:"userId" comment:"用户ID"` // 用户ID
	Username string `json:"username" comment:"用户名" vd:"len($)>0"`
	NickName string `json:"nickName" comment:"昵称" vd:"len($)>0"`
	Phone    string `json:"phone" comment:"手机号" vd:"len($)>0"`
	RoleIds  []int  `json:"roleIds" comment:"角色ID数组"`
	Avatar   string `json:"avatar" comment:"头像"`
	Sex      string `json:"sex" comment:"性别"`
	Email    string `json:"email" comment:"邮箱" vd:"len($)>0,email"`
	Remark   string `json:"remark" comment:"备注"`
	Status   string `json:"status" comment:"状态" default:"1"`
	common.ControlBy
}

func (s *SysUserUpdateReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Username = s.Username
	model.NickName = s.NickName
	model.Phone = s.Phone
	model.Avatar = s.Avatar
	model.Sex = s.Sex
	model.Email = s.Email
	model.Remark = s.Remark
	model.Status = s.Status

	// 拼接 RoleIds 字符串
	roleIdsStr := ""
	for i, roleId := range s.RoleIds {
		roleIdsStr += strconv.Itoa(roleId)
		if i < len(s.RoleIds)-1 {
			roleIdsStr += ","
		}
	}
	model.RoleIds = roleIdsStr
}

func (s *SysUserUpdateReq) GetId() interface{} {
	return s.UserId
}

type SysUserById struct {
	dto.ObjectById
	common.ControlBy
}

func (s *SysUserById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysUserById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysUser{}, nil
}

// PassWord 密码
type PassWord struct {
	NewPassword string `json:"newPassword" vd:"len($)>0"`
	OldPassword string `json:"oldPassword" vd:"len($)>0"`
}
