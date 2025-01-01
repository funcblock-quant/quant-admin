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

type BusTriangleOrdersRecord struct {
	api.Api
}

// GetPage 获取三角套利订单记录列表
// @Summary 获取三角套利订单记录列表
// @Description 获取三角套利订单记录列表
// @Tags 三角套利订单记录
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusTriangleOrdersRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-triangle-orders-record [get]
// @Security Bearer
func (e BusTriangleOrdersRecord) GetPage(c *gin.Context) {
	req := dto.BusTriangleOrdersRecordGetPageReq{}
	s := service.BusTriangleOrdersRecord{}
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
	list := make([]models.BusTriangleOrdersRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取三角套利订单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取三角套利订单记录
// @Summary 获取三角套利订单记录
// @Description 获取三角套利订单记录
// @Tags 三角套利订单记录
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusTriangleOrdersRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-triangle-orders-record/{id} [get]
// @Security Bearer
func (e BusTriangleOrdersRecord) Get(c *gin.Context) {
	req := dto.BusTriangleOrdersRecordGetReq{}
	s := service.BusTriangleOrdersRecord{}
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
	var object models.BusTriangleOrdersRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取三角套利订单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建三角套利订单记录
// @Summary 创建三角套利订单记录
// @Description 创建三角套利订单记录
// @Tags 三角套利订单记录
// @Accept application/json
// @Product application/json
// @Param data body dto.BusTriangleOrdersRecordInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-triangle-orders-record [post]
// @Security Bearer
func (e BusTriangleOrdersRecord) Insert(c *gin.Context) {
	req := dto.BusTriangleOrdersRecordInsertReq{}
	s := service.BusTriangleOrdersRecord{}
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
		e.Error(500, err, fmt.Sprintf("创建三角套利订单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改三角套利订单记录
// @Summary 修改三角套利订单记录
// @Description 修改三角套利订单记录
// @Tags 三角套利订单记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusTriangleOrdersRecordUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-triangle-orders-record/{id} [put]
// @Security Bearer
func (e BusTriangleOrdersRecord) Update(c *gin.Context) {
	req := dto.BusTriangleOrdersRecordUpdateReq{}
	s := service.BusTriangleOrdersRecord{}
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
		e.Error(500, err, fmt.Sprintf("修改三角套利订单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除三角套利订单记录
// @Summary 删除三角套利订单记录
// @Description 删除三角套利订单记录
// @Tags 三角套利订单记录
// @Param data body dto.BusTriangleOrdersRecordDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-triangle-orders-record [delete]
// @Security Bearer
func (e BusTriangleOrdersRecord) Delete(c *gin.Context) {
	s := service.BusTriangleOrdersRecord{}
	req := dto.BusTriangleOrdersRecordDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除三角套利订单记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
