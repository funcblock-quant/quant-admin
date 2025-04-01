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

type BusDexCexDepositWithdrawRecord struct {
	api.Api
}

// GetPage 获取DEX-CEX充提记录列表
// @Summary 获取DEX-CEX充提记录列表
// @Description 获取DEX-CEX充提记录列表
// @Tags DEX-CEX充提记录
// @Param orderType query string false "订单类型"
// @Param fromAddress query string false "from地址"
// @Param toAddress query string false "to地址"
// @Param cexUid query string false "cex账户uid"
// @Param cexOrderId query string false "交易所充提业务id"
// @Param txHash query string false "链上交易hash"
// @Param symbol query string false "充提币种"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusDexCexDepositWithdrawRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-cex-deposit-withdraw-record [get]
// @Security Bearer
func (e BusDexCexDepositWithdrawRecord) GetPage(c *gin.Context) {
    req := dto.BusDexCexDepositWithdrawRecordGetPageReq{}
    s := service.BusDexCexDepositWithdrawRecord{}
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
	list := make([]models.BusDexCexDepositWithdrawRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取DEX-CEX充提记录失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取DEX-CEX充提记录
// @Summary 获取DEX-CEX充提记录
// @Description 获取DEX-CEX充提记录
// @Tags DEX-CEX充提记录
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusDexCexDepositWithdrawRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-cex-deposit-withdraw-record/{id} [get]
// @Security Bearer
func (e BusDexCexDepositWithdrawRecord) Get(c *gin.Context) {
	req := dto.BusDexCexDepositWithdrawRecordGetReq{}
	s := service.BusDexCexDepositWithdrawRecord{}
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
	var object models.BusDexCexDepositWithdrawRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取DEX-CEX充提记录失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建DEX-CEX充提记录
// @Summary 创建DEX-CEX充提记录
// @Description 创建DEX-CEX充提记录
// @Tags DEX-CEX充提记录
// @Accept application/json
// @Product application/json
// @Param data body dto.BusDexCexDepositWithdrawRecordInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-dex-cex-deposit-withdraw-record [post]
// @Security Bearer
func (e BusDexCexDepositWithdrawRecord) Insert(c *gin.Context) {
    req := dto.BusDexCexDepositWithdrawRecordInsertReq{}
    s := service.BusDexCexDepositWithdrawRecord{}
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
		e.Error(500, err, fmt.Sprintf("创建DEX-CEX充提记录失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改DEX-CEX充提记录
// @Summary 修改DEX-CEX充提记录
// @Description 修改DEX-CEX充提记录
// @Tags DEX-CEX充提记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusDexCexDepositWithdrawRecordUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-dex-cex-deposit-withdraw-record/{id} [put]
// @Security Bearer
func (e BusDexCexDepositWithdrawRecord) Update(c *gin.Context) {
    req := dto.BusDexCexDepositWithdrawRecordUpdateReq{}
    s := service.BusDexCexDepositWithdrawRecord{}
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
		e.Error(500, err, fmt.Sprintf("修改DEX-CEX充提记录失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除DEX-CEX充提记录
// @Summary 删除DEX-CEX充提记录
// @Description 删除DEX-CEX充提记录
// @Tags DEX-CEX充提记录
// @Param data body dto.BusDexCexDepositWithdrawRecordDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-dex-cex-deposit-withdraw-record [delete]
// @Security Bearer
func (e BusDexCexDepositWithdrawRecord) Delete(c *gin.Context) {
    s := service.BusDexCexDepositWithdrawRecord{}
    req := dto.BusDexCexDepositWithdrawRecordDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除DEX-CEX充提记录失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
