package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
)

type StrategyDexCexTriangularArbitrageOpportunities struct {
	api.Api
}

// GetPage 获取CEX-DEX套利机会统计列表
// @Summary 获取CEX-DEX套利机会统计列表
// @Description 获取CEX-DEX套利机会统计列表
// @Tags CEX-DEX套利机会统计
// @Param instanceId query string false "实例id"
// @Param buyOnDex query string false "买方标识"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.StrategyDexCexTriangularArbitrageOpportunities}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-cex-triangular-arbitrage-opportunities [get]
// @Security Bearer
func (e StrategyDexCexTriangularArbitrageOpportunities) GetPage(c *gin.Context) {
	req := dto.StrategyDexCexTriangularArbitrageOpportunitiesGetPageReq{}
	s := service.StrategyDexCexTriangularArbitrageOpportunities{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.StrategyDexCexTriangularArbitrageOpportunities, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CEX-DEX套利机会统计失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CEX-DEX套利机会统计
// @Summary 获取CEX-DEX套利机会统计
// @Description 获取CEX-DEX套利机会统计
// @Tags CEX-DEX套利机会统计
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.StrategyDexCexTriangularArbitrageOpportunities} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-cex-triangular-arbitrage-opportunities/{id} [get]
// @Security Bearer
func (e StrategyDexCexTriangularArbitrageOpportunities) Get(c *gin.Context) {
	req := dto.StrategyDexCexTriangularArbitrageOpportunitiesGetReq{}
	s := service.StrategyDexCexTriangularArbitrageOpportunities{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.StrategyDexCexTriangularArbitrageOpportunities

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CEX-DEX套利机会统计失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建CEX-DEX套利机会统计
// @Summary 创建CEX-DEX套利机会统计
// @Description 创建CEX-DEX套利机会统计
// @Tags CEX-DEX套利机会统计
// @Accept application/json
// @Product application/json
// @Param data body dto.StrategyDexCexTriangularArbitrageOpportunitiesInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-dex-cex-triangular-arbitrage-opportunities [post]
// @Security Bearer
func (e StrategyDexCexTriangularArbitrageOpportunities) Insert(c *gin.Context) {
	req := dto.StrategyDexCexTriangularArbitrageOpportunitiesInsertReq{}
	s := service.StrategyDexCexTriangularArbitrageOpportunities{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建CEX-DEX套利机会统计失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CEX-DEX套利机会统计
// @Summary 修改CEX-DEX套利机会统计
// @Description 修改CEX-DEX套利机会统计
// @Tags CEX-DEX套利机会统计
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.StrategyDexCexTriangularArbitrageOpportunitiesUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-dex-cex-triangular-arbitrage-opportunities/{id} [put]
// @Security Bearer
func (e StrategyDexCexTriangularArbitrageOpportunities) Update(c *gin.Context) {
	req := dto.StrategyDexCexTriangularArbitrageOpportunitiesUpdateReq{}
	s := service.StrategyDexCexTriangularArbitrageOpportunities{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CEX-DEX套利机会统计失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除CEX-DEX套利机会统计
// @Summary 删除CEX-DEX套利机会统计
// @Description 删除CEX-DEX套利机会统计
// @Tags CEX-DEX套利机会统计
// @Param data body dto.StrategyDexCexTriangularArbitrageOpportunitiesDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-dex-cex-triangular-arbitrage-opportunities [delete]
// @Security Bearer
func (e StrategyDexCexTriangularArbitrageOpportunities) Delete(c *gin.Context) {
	s := service.StrategyDexCexTriangularArbitrageOpportunities{}
	req := dto.StrategyDexCexTriangularArbitrageOpportunitiesDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除CEX-DEX套利机会统计失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
