package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"quanta-admin/app/business/apis"
	"quanta-admin/common/actions"
	"quanta-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerBusStrategyInstanceRouter)
}

// registerBusStrategyInstanceRouter
func registerBusStrategyInstanceRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.BusStrategyInstance{}
	r := v1.Group("/strategy-instance").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
	}

	r = v1.PUT("startStrategyInstance/:id", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.StartInstance)
	r = v1.PUT("stopStrategyInstance/:id", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.StopInstance)
	r = v1.POST("batchStartStrategyInstance", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.BatchStartInstance)
	r = v1.POST("batchStopStrategyInstance", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.BatchStopInstance)

}
