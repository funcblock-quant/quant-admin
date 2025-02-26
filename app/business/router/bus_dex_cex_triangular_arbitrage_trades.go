package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"quanta-admin/app/business/apis"
	"quanta-admin/common/actions"
	"quanta-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerStrategyDexCexTriangularArbitrageTradesRouter)
}

// registerStrategyDexCexTriangularArbitrageTradesRouter
func registerStrategyDexCexTriangularArbitrageTradesRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.StrategyDexCexTriangularArbitrageTrades{}
	r := v1.Group("/bus-dex-cex-triangular-arbitrage-trades").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
	}
	v1.GET("/getDexCexTriangularTraderStatistics", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.GetDexCexTriangularTraderStatistics)
}
