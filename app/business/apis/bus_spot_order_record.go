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

type BusSpotOrderRecord struct {
	api.Api
}

// GetPage 获取现货交易记录列表
// @Summary 获取现货交易记录列表
// @Description 获取现货交易记录列表
// @Tags 现货交易记录
// @Param arbitrageId query string false "套利记录id"
// @Param side query string false "买卖方向"
// @Param orderId query string false "交易所订单id"
// @Param orderClientId query string false "策略端id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusSpotOrderRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/spot-order [get]
// @Security Bearer
func (e BusSpotOrderRecord) GetPage(c *gin.Context) {
	req := dto.BusSpotOrderRecordGetPageReq{}
	s := service.BusSpotOrderRecord{}
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
	list := make([]models.BusSpotOrderRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取现货交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取现货交易记录
// @Summary 获取现货交易记录
// @Description 获取现货交易记录
// @Tags 现货交易记录
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusSpotOrderRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/spot-order/{id} [get]
// @Security Bearer
func (e BusSpotOrderRecord) Get(c *gin.Context) {
	req := dto.BusSpotOrderRecordGetReq{}
	s := service.BusSpotOrderRecord{}
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
	var object models.BusSpotOrderRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取现货交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建现货交易记录
// @Summary 创建现货交易记录
// @Description 创建现货交易记录
// @Tags 现货交易记录
// @Accept application/json
// @Product application/json
// @Param data body dto.BusSpotOrderRecordInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/spot-order [post]
// @Security Bearer
func (e BusSpotOrderRecord) Insert(c *gin.Context) {
	req := dto.BusSpotOrderRecordInsertReq{}
	s := service.BusSpotOrderRecord{}
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
		e.Error(500, err, fmt.Sprintf("创建现货交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改现货交易记录
// @Summary 修改现货交易记录
// @Description 修改现货交易记录
// @Tags 现货交易记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusSpotOrderRecordUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/spot-order/{id} [put]
// @Security Bearer
func (e BusSpotOrderRecord) Update(c *gin.Context) {
	req := dto.BusSpotOrderRecordUpdateReq{}
	s := service.BusSpotOrderRecord{}
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
		e.Error(500, err, fmt.Sprintf("修改现货交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除现货交易记录
// @Summary 删除现货交易记录
// @Description 删除现货交易记录
// @Tags 现货交易记录
// @Param data body dto.BusSpotOrderRecordDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/spot-order [delete]
// @Security Bearer
func (e BusSpotOrderRecord) Delete(c *gin.Context) {
	s := service.BusSpotOrderRecord{}
	req := dto.BusSpotOrderRecordDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除现货交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
