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

type BusExchangeAccountInfo struct {
	api.Api
}

// GetPage 获取账户配置列表
// @Summary 获取账户配置列表
// @Description 获取账户配置列表
// @Tags 账户配置
// @Param accountName query string false "钱包名称"
// @Param platformId query string false "id"
// @Param platformName query string false "交易所名称"
// @Param accountType query string false "账户类型"
// @Param status query string false "状态"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusExchangeAccountInfo}} "{"code": 200, "data": [...]}"
// @Router /api/v1/exchange-account [get]
// @Security Bearer
func (e BusExchangeAccountInfo) GetPage(c *gin.Context) {
	req := dto.BusExchangeAccountInfoGetPageReq{}
	s := service.BusExchangeAccountInfo{}
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
	list := make([]models.BusExchangeAccountInfo, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取账户配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取账户配置
// @Summary 获取账户配置
// @Description 获取账户配置
// @Tags 账户配置
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusExchangeAccountInfo} "{"code": 200, "data": [...]}"
// @Router /api/v1/exchange-account/{id} [get]
// @Security Bearer
func (e BusExchangeAccountInfo) Get(c *gin.Context) {
	req := dto.BusExchangeAccountInfoGetReq{}
	s := service.BusExchangeAccountInfo{}
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
	var object models.BusExchangeAccountInfo

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取账户配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建账户配置
// @Summary 创建账户配置
// @Description 创建账户配置
// @Tags 账户配置
// @Accept application/json
// @Product application/json
// @Param data body dto.BusExchangeAccountInfoInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/exchange-account [post]
// @Security Bearer
func (e BusExchangeAccountInfo) Insert(c *gin.Context) {
	req := dto.BusExchangeAccountInfoInsertReq{}
	s := service.BusExchangeAccountInfo{}
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
		e.Error(500, err, fmt.Sprintf("创建账户配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改账户配置
// @Summary 修改账户配置
// @Description 修改账户配置
// @Tags 账户配置
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusExchangeAccountInfoUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/exchange-account/{id} [put]
// @Security Bearer
func (e BusExchangeAccountInfo) Update(c *gin.Context) {
	req := dto.BusExchangeAccountInfoUpdateReq{}
	s := service.BusExchangeAccountInfo{}
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
		e.Error(500, err, fmt.Sprintf("修改账户配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除账户配置
// @Summary 删除账户配置
// @Description 删除账户配置
// @Tags 账户配置
// @Param data body dto.BusExchangeAccountInfoDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/exchange-account [delete]
// @Security Bearer
func (e BusExchangeAccountInfo) Delete(c *gin.Context) {
	s := service.BusExchangeAccountInfo{}
	req := dto.BusExchangeAccountInfoDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除账户配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

func (e BusExchangeAccountInfo) GetAccountListByGroupId(c *gin.Context) {
	s := service.BusExchangeAccountInfo{}
	req := dto.BusGroupAccountInfoGetReq{}
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

	list := make([]models.BusExchangeAccountInfo, 0)
	err = s.GetAccountListByGroupId(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除账户配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "删除成功")
}
