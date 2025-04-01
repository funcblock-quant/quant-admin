package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"quanta-admin/app/business/apis"
	"quanta-admin/common/actions"
	"quanta-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerBusExchangeAccountInfoRouter)
}

// registerBusExchangeAccountInfoRouter
func registerBusExchangeAccountInfoRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.BusExchangeAccountInfo{}
	r := v1.Group("/exchange-account").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
	}
	r = v1.GET("/queryAccountByGroupId/:groupId", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.GetAccountListByGroupId)
	r = v1.GET("/queryExchangeListInUse", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.QueryExchangeListInUse)
	r = v1.POST("/getPortfolioUnwindingInfo", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.GetPortfolioUnwindingInfo)
	r = v1.POST("/portfolioUnwinding", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.PortfolioUnwinding)
}
