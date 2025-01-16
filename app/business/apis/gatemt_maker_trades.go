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

type GatemtMakerTrades struct {
	api.Api
}

// GetPage 获取GateMT做市交易明细列表
// @Summary 获取GateMT做市交易明细列表
// @Description 获取GateMT做市交易明细列表
// @Tags GateMT做市交易明细
// @Param clientOrderId query string false ""
// @Param binanceLimitOrderId query string false ""
// @Param binanceMarketOrderId query string false ""
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GatemtMakerTrades}} "{"code": 200, "data": [...]}"
// @Router /api/v1/market-making-trades [get]
// @Security Bearer
func (e GatemtMakerTrades) GetPage(c *gin.Context) {
	req := dto.GatemtMakerTradesGetPageReq{}
	s := service.GatemtMakerTrades{}
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
	list := make([]models.GatemtMakerTrades, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GateMT做市交易明细失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// QueryTradesWithClientOrderId 根据clientOrderId获取GateMT做市交易明细列表
// @Summary 根据clientOrderId获取GateMT做市交易明细列表
// @Description 根据clientOrderId获取GateMT做市交易明细列表
// @Tags GateMT做市交易明细
// @Param clientOrderId query string false ""
// @Param binanceLimitOrderId query string false ""
// @Param binanceMarketOrderId query string false ""
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GatemtMakerTrades}} "{"code": 200, "data": [...]}"
// @Router /api/v1/queryMarketMakingTradesByOrderId/:orderId [get]
// @Security Bearer
func (e GatemtMakerTrades) QueryTradesWithClientOrderId(c *gin.Context) {
	req := dto.GatemtMakerTradesGetListReq{}
	s := service.GatemtMakerTrades{}
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
	list := make([]models.GatemtMakerTrades, 0)

	err = s.QueryTradesWithClientOrderId(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GateMT做市交易明细失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}

// Get 获取GateMT做市交易明细
// @Summary 获取GateMT做市交易明细
// @Description 获取GateMT做市交易明细
// @Tags GateMT做市交易明细
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GatemtMakerTrades} "{"code": 200, "data": [...]}"
// @Router /api/v1/market-making-trades/{id} [get]
// @Security Bearer
func (e GatemtMakerTrades) Get(c *gin.Context) {
	req := dto.GatemtMakerTradesGetReq{}
	s := service.GatemtMakerTrades{}
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
	var object models.GatemtMakerTrades

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GateMT做市交易明细失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建GateMT做市交易明细
// @Summary 创建GateMT做市交易明细
// @Description 创建GateMT做市交易明细
// @Tags GateMT做市交易明细
// @Accept application/json
// @Product application/json
// @Param data body dto.GatemtMakerTradesInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/market-making-trades [post]
// @Security Bearer
func (e GatemtMakerTrades) Insert(c *gin.Context) {
	req := dto.GatemtMakerTradesInsertReq{}
	s := service.GatemtMakerTrades{}
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
		e.Error(500, err, fmt.Sprintf("创建GateMT做市交易明细失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GateMT做市交易明细
// @Summary 修改GateMT做市交易明细
// @Description 修改GateMT做市交易明细
// @Tags GateMT做市交易明细
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GatemtMakerTradesUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/market-making-trades/{id} [put]
// @Security Bearer
func (e GatemtMakerTrades) Update(c *gin.Context) {
	req := dto.GatemtMakerTradesUpdateReq{}
	s := service.GatemtMakerTrades{}
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
		e.Error(500, err, fmt.Sprintf("修改GateMT做市交易明细失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除GateMT做市交易明细
// @Summary 删除GateMT做市交易明细
// @Description 删除GateMT做市交易明细
// @Tags GateMT做市交易明细
// @Param data body dto.GatemtMakerTradesDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/market-making-trades [delete]
// @Security Bearer
func (e GatemtMakerTrades) Delete(c *gin.Context) {
	s := service.GatemtMakerTrades{}
	req := dto.GatemtMakerTradesDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除GateMT做市交易明细失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
