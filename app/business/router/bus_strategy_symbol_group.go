package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"quanta-admin/app/business/apis"
	"quanta-admin/common/actions"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerBusStrategySymbolGroupRouter)
}

// registerBusStrategySymbolGroupRouter
func registerBusStrategySymbolGroupRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.BusStrategySymbolGroup{}
	r := v1.Group("/symbol-group").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
	}
}
