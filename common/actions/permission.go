package actions

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"gorm.io/gorm"
)

type DataPermission struct {
	DataScope string
	UserId    int
	RoleIds   []int
}

func PermissionAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		var p = new(DataPermission)
		if userId := user.GetUserIdStr(c); userId != "" {
			p, err = newDataPermission(db, userId)
			if err != nil {
				log.Errorf("MsgID[%s] PermissionAction error: %s", msgID, err)
				response.Error(c, 500, err, "权限范围鉴定错误")
				c.Abort()
				return
			}
		}
		c.Set(PermissionKey, p)
		c.Next()
	}
}

func newDataPermission(tx *gorm.DB, userId interface{}) (*DataPermission, error) {
	var err error
	p := &DataPermission{}
	var roleIdsStr string

	err = tx.Table("sys_user").
		Select("sys_user.user_id", "sys_user.role_ids", "sys_role.data_scope").
		Joins("left join sys_role on FIND_IN_SET(sys_role.role_id, sys_user.role_ids) > 0"). // 使用 FIND_IN_SET
		Where("sys_user.user_id = ?", userId).
		Scan(&struct {
			UserId    int
			RoleIds   string
			DataScope string
		}{
			UserId:    p.UserId,
			RoleIds:   roleIdsStr,
			DataScope: p.DataScope,
		}).Error
	if err != nil {
		err = errors.New("获取用户数据出错 msg:" + err.Error())
		return nil, err
	}

	// 解析 role_ids 字符串为 []int
	if roleIdsStr != "" {
		roleIds := strings.Split(roleIdsStr, ",")
		for _, roleIdStr := range roleIds {
			roleId, err := strconv.Atoi(roleIdStr)
			if err != nil {
				return nil, err // 如果转换失败，返回错误
			}
			p.RoleIds = append(p.RoleIds, roleId)
		}
	}
	return p, nil
}

func Permission(tableName string, p *DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !config.ApplicationConfig.EnableDP {
			return db
		}
		switch p.DataScope {
		case "2":
			return db.Where(tableName+".create_by = ?", p.UserId)
		default:
			return db
		}
	}
}

func getPermissionFromContext(c *gin.Context) *DataPermission {
	p := new(DataPermission)
	if pm, ok := c.Get(PermissionKey); ok {
		switch pm.(type) {
		case *DataPermission:
			p = pm.(*DataPermission)
		}
	}
	return p
}

// GetPermissionFromContext 提供非action写法数据范围约束
func GetPermissionFromContext(c *gin.Context) *DataPermission {
	return getPermissionFromContext(c)
}
