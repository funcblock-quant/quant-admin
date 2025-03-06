package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"quanta-admin/app/business/apis"
	"quanta-admin/common/actions"
	"quanta-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerBusDexCexTriangularObserverRouter)
}

// registerBusDexCexTriangularObserverRouter
func registerBusDexCexTriangularObserverRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.BusDexCexTriangularObserver{}
	r := v1.Group("/dex-cex-triangular-observer").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
	}
	v1.POST("/batchAddBusDexCexTriangularObserver", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.BatchInsert)
	v1.GET("/busDexCexTriangularSymbolList", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.GetSymbolList)
	v1.POST("/busDexCexTriangularStartTrader", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.StartTrader)
	v1.POST("/busDexCexTriangularStopTrader", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.StopTrader)
	v1.PUT("/busDexCexTriangularUpdateObserver", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.UpdateObserver)
	v1.PUT("/busDexCexTriangularUpdateTrader", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.UpdateTrader)
	v1.PUT("/busDexCexTriangularUpdateWaterLevel", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.UpdateWaterLevel)
	v1.GET("/busDexCexTriangularGetGlobalWaterLevel", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.GetGlobalWaterLevelState)
	v1.POST("/busDexCexTriangularUpdateGlobalWaterLevel", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.UpdateGlobalWaterLevel)
	v1.GET("/busDexCexTriangularGetRiskConfig", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.GetGlobalRiskConfigState)
	v1.POST("/busDexCexTriangularUpdateRiskConfig", authMiddleware.MiddlewareFunc(), middleware.AuthCheckRole(), actions.PermissionAction(), api.UpdateGlobalRiskConfig)
}
