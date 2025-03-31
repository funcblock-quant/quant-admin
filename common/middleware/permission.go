package middleware

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2/util"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

// AuthCheckRole 权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		data, _ := c.Get(jwtauth.JwtPayloadKey)
		v := data.(jwtauth.MapClaims)
		e := sdk.Runtime.GetCasbinKey(c.Request.Host)
		var res, casbinExclude bool
		var err error
		var hasPermission bool // 添加一个变量，用于标记是否有权限
		var enforceErr error   // 添加一个变量，用于存储 Enforce 错误
		//检查权限
		roleKeyList := strings.Split(v["rolekey"].(string), ",")
		for _, roleKey := range roleKeyList {
			if roleKey == "admin" {
				res = true
				c.Next()
				return
			}
		}

		for _, i := range CasbinExclude {
			if util.KeyMatch2(c.Request.URL.Path, i.Url) && c.Request.Method == i.Method {
				casbinExclude = true
				break
			}
		}
		if casbinExclude {
			log.Infof("Casbin exclusion, no validation method:%s path:%s", c.Request.Method, c.Request.URL.Path)
			c.Next()
			return
		}
		// Casbin 权限验证
		for _, roleKey := range roleKeyList {
			res, err = e.Enforce(roleKey, c.Request.URL.Path, c.Request.Method)
			if err != nil && enforceErr == nil { // 记录第一个 Enforce 错误
				enforceErr = err
			}
			if res {
				hasPermission = true // 标记有权限
				break
			}
		}

		if enforceErr != nil { // 检查是否有 Enforce 错误
			log.Errorf("AuthCheckRole error:%s method:%s path:%s", enforceErr, c.Request.Method, c.Request.URL.Path)
			response.Error(c, 500, enforceErr, "")
			c.Abort()
			return
		}

		if hasPermission { // 检查是否有权限
			log.Infof("isTrue: %v role: %s method: %s path: %s", res, v["rolekey"], c.Request.Method, c.Request.URL.Path)
			c.Next()
		} else {
			log.Warnf("isTrue: %v role: %s method: %s path: %s message: %s", res, v["rolekey"], c.Request.Method, c.Request.URL.Path, "当前request无权限，请管理员确认！")
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "对不起，您没有该接口访问权限，请联系管理员",
			})
			c.Abort()
			return
		}

	}
}
