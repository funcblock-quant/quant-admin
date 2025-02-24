package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"quanta-admin/app/business/apis"
	"quanta-admin/common/actions"
	"quanta-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerBusPriceTriggerStrategyInstanceRouter)
}

// registerBusPriceTriggerStrategyInstanceRouter
func registerBusPriceTriggerStrategyInstanceRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.BusPriceTriggerStrategyInstance{}
	r := v1.Group("/bus-price-trigger-strategy").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
	}
	r = v1.POST("/stopTriggerInstance", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.StopInstance)
	r = v1.GET("/getSymbolList", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.GetSymbolList)
	r = v1.PUT("/updateProfitTarget", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.UpdateProfitTarget)
	r = v1.PUT("/updateExecuteNum", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.UpdateExecuteNum)

}
