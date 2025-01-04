package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"quanta-admin/app/business/apis"
	"quanta-admin/common/actions"
	"quanta-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerBusDexCexTriangularArbitrageRecordRouter)
}

// registerBusDexCexTriangularArbitrageRecordRouter
func registerBusDexCexTriangularArbitrageRecordRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.BusDexCexTriangularArbitrageRecord{}
	r := v1.Group("/dex-cex-triangular").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
	}
	v1.POST("/market/dex-cex-arbitrage-chance", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.GetArbitrageOpportunity)
}
