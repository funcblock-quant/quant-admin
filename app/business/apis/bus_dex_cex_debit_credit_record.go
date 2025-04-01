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

type BusDexCexDebitCreditRecord struct {
	api.Api
}

// GetPage 获取CEX借贷记录列表
// @Summary 获取CEX借贷记录列表
// @Description 获取CEX借贷记录列表
// @Tags CEX借贷记录
// @Param accountName query string false "交易所账户名称"
// @Param uid query string false "交易所账户uid"
// @Param exchangeType query string false "交易所"
// @Param debitType query string false "类型"
// @Param status query string false "状态"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusDexCexDebitCreditRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-cex-debit-credit-record [get]
// @Security Bearer
func (e BusDexCexDebitCreditRecord) GetPage(c *gin.Context) {
    req := dto.BusDexCexDebitCreditRecordGetPageReq{}
    s := service.BusDexCexDebitCreditRecord{}
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
	list := make([]models.BusDexCexDebitCreditRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CEX借贷记录失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CEX借贷记录
// @Summary 获取CEX借贷记录
// @Description 获取CEX借贷记录
// @Tags CEX借贷记录
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusDexCexDebitCreditRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-cex-debit-credit-record/{id} [get]
// @Security Bearer
func (e BusDexCexDebitCreditRecord) Get(c *gin.Context) {
	req := dto.BusDexCexDebitCreditRecordGetReq{}
	s := service.BusDexCexDebitCreditRecord{}
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
	var object models.BusDexCexDebitCreditRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CEX借贷记录失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建CEX借贷记录
// @Summary 创建CEX借贷记录
// @Description 创建CEX借贷记录
// @Tags CEX借贷记录
// @Accept application/json
// @Product application/json
// @Param data body dto.BusDexCexDebitCreditRecordInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-dex-cex-debit-credit-record [post]
// @Security Bearer
func (e BusDexCexDebitCreditRecord) Insert(c *gin.Context) {
    req := dto.BusDexCexDebitCreditRecordInsertReq{}
    s := service.BusDexCexDebitCreditRecord{}
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
		e.Error(500, err, fmt.Sprintf("创建CEX借贷记录失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CEX借贷记录
// @Summary 修改CEX借贷记录
// @Description 修改CEX借贷记录
// @Tags CEX借贷记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusDexCexDebitCreditRecordUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-dex-cex-debit-credit-record/{id} [put]
// @Security Bearer
func (e BusDexCexDebitCreditRecord) Update(c *gin.Context) {
    req := dto.BusDexCexDebitCreditRecordUpdateReq{}
    s := service.BusDexCexDebitCreditRecord{}
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
		e.Error(500, err, fmt.Sprintf("修改CEX借贷记录失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除CEX借贷记录
// @Summary 删除CEX借贷记录
// @Description 删除CEX借贷记录
// @Tags CEX借贷记录
// @Param data body dto.BusDexCexDebitCreditRecordDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-dex-cex-debit-credit-record [delete]
// @Security Bearer
func (e BusDexCexDebitCreditRecord) Delete(c *gin.Context) {
    s := service.BusDexCexDebitCreditRecord{}
    req := dto.BusDexCexDebitCreditRecordDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除CEX借贷记录失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
