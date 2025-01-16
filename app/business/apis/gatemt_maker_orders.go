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

type GatemtMakerOrders struct {
	api.Api
}

// GetPage 获取GateMT做市列表
// @Summary 获取GateMT做市列表
// @Description 获取GateMT做市列表
// @Tags GateMT做市
// @Param clientOrderId query string false "Client order id"
// @Param symbol query string false ""
// @Param exchangeOrderId query string false ""
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GatemtMakerOrders}} "{"code": 200, "data": [...]}"
// @Router /api/v1/market-making-record [get]
// @Security Bearer
func (e GatemtMakerOrders) GetPage(c *gin.Context) {
	req := dto.GatemtMakerOrdersGetPageReq{}
	s := service.GatemtMakerOrders{}
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
	list := make([]models.GatemtMakerOrders, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GateMT做市失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GateMT做市
// @Summary 获取GateMT做市
// @Description 获取GateMT做市
// @Tags GateMT做市
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GatemtMakerOrders} "{"code": 200, "data": [...]}"
// @Router /api/v1/market-making-record/{id} [get]
// @Security Bearer
func (e GatemtMakerOrders) Get(c *gin.Context) {
	req := dto.GatemtMakerOrdersGetReq{}
	s := service.GatemtMakerOrders{}
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
	var object models.GatemtMakerOrders

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GateMT做市失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建GateMT做市
// @Summary 创建GateMT做市
// @Description 创建GateMT做市
// @Tags GateMT做市
// @Accept application/json
// @Product application/json
// @Param data body dto.GatemtMakerOrdersInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/market-making-record [post]
// @Security Bearer
func (e GatemtMakerOrders) Insert(c *gin.Context) {
	req := dto.GatemtMakerOrdersInsertReq{}
	s := service.GatemtMakerOrders{}
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
		e.Error(500, err, fmt.Sprintf("创建GateMT做市失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GateMT做市
// @Summary 修改GateMT做市
// @Description 修改GateMT做市
// @Tags GateMT做市
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GatemtMakerOrdersUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/market-making-record/{id} [put]
// @Security Bearer
func (e GatemtMakerOrders) Update(c *gin.Context) {
	req := dto.GatemtMakerOrdersUpdateReq{}
	s := service.GatemtMakerOrders{}
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
		e.Error(500, err, fmt.Sprintf("修改GateMT做市失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除GateMT做市
// @Summary 删除GateMT做市
// @Description 删除GateMT做市
// @Tags GateMT做市
// @Param data body dto.GatemtMakerOrdersDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/market-making-record [delete]
// @Security Bearer
func (e GatemtMakerOrders) Delete(c *gin.Context) {
	s := service.GatemtMakerOrders{}
	req := dto.GatemtMakerOrdersDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除GateMT做市失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
