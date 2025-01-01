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

type BusFuturePositionRecord struct {
	api.Api
}

// GetPage 获取合约交易记录列表
// @Summary 获取合约交易记录列表
// @Description 获取合约交易记录列表
// @Tags 合约交易记录
// @Param arbitrageId query string false "套利记录id"
// @Param symbol query string false "交易币种"
// @Param orderId query string false "交易所订单id"
// @Param orderClientId query string false "策略端id"
// @Param positionSide query string false "持仓方向"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusFuturePositionRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/future-order [get]
// @Security Bearer
func (e BusFuturePositionRecord) GetPage(c *gin.Context) {
	req := dto.BusFuturePositionRecordGetPageReq{}
	s := service.BusFuturePositionRecord{}
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
	list := make([]models.BusFuturePositionRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取合约交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取合约交易记录
// @Summary 获取合约交易记录
// @Description 获取合约交易记录
// @Tags 合约交易记录
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusFuturePositionRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/future-order/{id} [get]
// @Security Bearer
func (e BusFuturePositionRecord) Get(c *gin.Context) {
	req := dto.BusFuturePositionRecordGetReq{}
	s := service.BusFuturePositionRecord{}
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
	var object models.BusFuturePositionRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取合约交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建合约交易记录
// @Summary 创建合约交易记录
// @Description 创建合约交易记录
// @Tags 合约交易记录
// @Accept application/json
// @Product application/json
// @Param data body dto.BusFuturePositionRecordInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/future-order [post]
// @Security Bearer
func (e BusFuturePositionRecord) Insert(c *gin.Context) {
	req := dto.BusFuturePositionRecordInsertReq{}
	s := service.BusFuturePositionRecord{}
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
		e.Error(500, err, fmt.Sprintf("创建合约交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改合约交易记录
// @Summary 修改合约交易记录
// @Description 修改合约交易记录
// @Tags 合约交易记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusFuturePositionRecordUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/future-order/{id} [put]
// @Security Bearer
func (e BusFuturePositionRecord) Update(c *gin.Context) {
	req := dto.BusFuturePositionRecordUpdateReq{}
	s := service.BusFuturePositionRecord{}
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
		e.Error(500, err, fmt.Sprintf("修改合约交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除合约交易记录
// @Summary 删除合约交易记录
// @Description 删除合约交易记录
// @Tags 合约交易记录
// @Param data body dto.BusFuturePositionRecordDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/future-order [delete]
// @Security Bearer
func (e BusFuturePositionRecord) Delete(c *gin.Context) {
	s := service.BusFuturePositionRecord{}
	req := dto.BusFuturePositionRecordDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除合约交易记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
