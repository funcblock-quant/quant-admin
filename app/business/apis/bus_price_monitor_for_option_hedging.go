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

type BusPriceMonitorForOptionHedging struct {
	api.Api
}

// GetPage 获取价格触发下单记录列表
// @Summary 获取价格触发下单记录列表
// @Description 获取价格触发下单记录列表
// @Tags 价格触发下单记录
// @Param monitoredOpenedNum query string false "监控的开单数量"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusPriceMonitorForOptionHedging}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-price-monitor [get]
// @Security Bearer
func (e BusPriceMonitorForOptionHedging) GetPage(c *gin.Context) {
	req := dto.BusPriceMonitorForOptionHedgingGetPageReq{}
	s := service.BusPriceMonitorForOptionHedging{}
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
	list := make([]models.BusPriceMonitorForOptionHedging, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取价格触发下单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取价格触发下单记录
// @Summary 获取价格触发下单记录
// @Description 获取价格触发下单记录
// @Tags 价格触发下单记录
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusPriceMonitorForOptionHedging} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-price-monitor/{id} [get]
// @Security Bearer
func (e BusPriceMonitorForOptionHedging) Get(c *gin.Context) {
	req := dto.BusPriceMonitorForOptionHedgingGetReq{}
	s := service.BusPriceMonitorForOptionHedging{}
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
	var object models.BusPriceMonitorForOptionHedging

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取价格触发下单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建价格触发下单记录
// @Summary 创建价格触发下单记录
// @Description 创建价格触发下单记录
// @Tags 价格触发下单记录
// @Accept application/json
// @Product application/json
// @Param data body dto.BusPriceMonitorForOptionHedgingInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-price-monitor [post]
// @Security Bearer
func (e BusPriceMonitorForOptionHedging) Insert(c *gin.Context) {
	req := dto.BusPriceMonitorForOptionHedgingInsertReq{}
	s := service.BusPriceMonitorForOptionHedging{}
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
		e.Error(500, err, fmt.Sprintf("创建价格触发下单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改价格触发下单记录
// @Summary 修改价格触发下单记录
// @Description 修改价格触发下单记录
// @Tags 价格触发下单记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusPriceMonitorForOptionHedgingUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-price-monitor/{id} [put]
// @Security Bearer
func (e BusPriceMonitorForOptionHedging) Update(c *gin.Context) {
	req := dto.BusPriceMonitorForOptionHedgingUpdateReq{}
	s := service.BusPriceMonitorForOptionHedging{}
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
		e.Error(500, err, fmt.Sprintf("修改价格触发下单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除价格触发下单记录
// @Summary 删除价格触发下单记录
// @Description 删除价格触发下单记录
// @Tags 价格触发下单记录
// @Param data body dto.BusPriceMonitorForOptionHedgingDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-price-monitor [delete]
// @Security Bearer
func (e BusPriceMonitorForOptionHedging) Delete(c *gin.Context) {
	s := service.BusPriceMonitorForOptionHedging{}
	req := dto.BusPriceMonitorForOptionHedgingDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除价格触发下单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
