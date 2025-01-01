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

type BusStrategyExchange struct {
	api.Api
}

// GetPage 获取交易所配置列表
// @Summary 获取交易所配置列表
// @Description 获取交易所配置列表
// @Tags 交易所配置
// @Param exchangeName query string false "名称"
// @Param exchangeType query string false "交易所类型"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusStrategyExchange}} "{"code": 200, "data": [...]}"
// @Router /api/v1/exchange [get]
// @Security Bearer
func (e BusStrategyExchange) GetPage(c *gin.Context) {
	req := dto.BusStrategyExchangeGetPageReq{}
	s := service.BusStrategyExchange{}
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
	list := make([]models.BusStrategyExchange, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取交易所配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取交易所配置
// @Summary 获取交易所配置
// @Description 获取交易所配置
// @Tags 交易所配置
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusStrategyExchange} "{"code": 200, "data": [...]}"
// @Router /api/v1/exchange/{id} [get]
// @Security Bearer
func (e BusStrategyExchange) Get(c *gin.Context) {
	req := dto.BusStrategyExchangeGetReq{}
	s := service.BusStrategyExchange{}
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
	var object models.BusStrategyExchange

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取交易所配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建交易所配置
// @Summary 创建交易所配置
// @Description 创建交易所配置
// @Tags 交易所配置
// @Accept application/json
// @Product application/json
// @Param data body dto.BusStrategyExchangeInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/exchange [post]
// @Security Bearer
func (e BusStrategyExchange) Insert(c *gin.Context) {
	req := dto.BusStrategyExchangeInsertReq{}
	s := service.BusStrategyExchange{}
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
		e.Error(500, err, fmt.Sprintf("创建交易所配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改交易所配置
// @Summary 修改交易所配置
// @Description 修改交易所配置
// @Tags 交易所配置
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusStrategyExchangeUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/exchange/{id} [put]
// @Security Bearer
func (e BusStrategyExchange) Update(c *gin.Context) {
	req := dto.BusStrategyExchangeUpdateReq{}
	s := service.BusStrategyExchange{}
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
		e.Error(500, err, fmt.Sprintf("修改交易所配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除交易所配置
// @Summary 删除交易所配置
// @Description 删除交易所配置
// @Tags 交易所配置
// @Param data body dto.BusStrategyExchangeDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/exchange [delete]
// @Security Bearer
func (e BusStrategyExchange) Delete(c *gin.Context) {
	s := service.BusStrategyExchange{}
	req := dto.BusStrategyExchangeDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除交易所配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
