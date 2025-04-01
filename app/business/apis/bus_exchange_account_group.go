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

type BusExchangeAccountGroup struct {
	api.Api
}

// GetPage 获取账号组设置列表
// @Summary 获取账号组设置列表
// @Description 获取账号组设置列表
// @Tags 账号组设置
// @Param groupName query string false "交易所账户组"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusExchangeAccountGroup}} "{"code": 200, "data": [...]}"
// @Router /api/v1/account-group [get]
// @Security Bearer
func (e BusExchangeAccountGroup) GetPage(c *gin.Context) {
	req := dto.BusExchangeAccountGroupGetPageReq{}
	s := service.BusExchangeAccountGroup{}
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
	list := make([]models.BusExchangeAccountGroup, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取账号组设置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取账号组设置
// @Summary 获取账号组设置
// @Description 获取账号组设置
// @Tags 账号组设置
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusExchangeAccountGroup} "{"code": 200, "data": [...]}"
// @Router /api/v1/account-group/{id} [get]
// @Security Bearer
func (e BusExchangeAccountGroup) Get(c *gin.Context) {
	req := dto.BusExchangeAccountGroupGetReq{}
	s := service.BusExchangeAccountGroup{}
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
	var object models.BusExchangeAccountGroup

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取账号组设置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建账号组设置
// @Summary 创建账号组设置
// @Description 创建账号组设置
// @Tags 账号组设置
// @Accept application/json
// @Product application/json
// @Param data body dto.BusExchangeAccountGroupInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/account-group [post]
// @Security Bearer
func (e BusExchangeAccountGroup) Insert(c *gin.Context) {
	req := dto.BusExchangeAccountGroupInsertReq{}
	s := service.BusExchangeAccountGroup{}
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
		e.Error(500, err, fmt.Sprintf("创建账号组设置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改账号组设置
// @Summary 修改账号组设置
// @Description 修改账号组设置
// @Tags 账号组设置
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusExchangeAccountGroupUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/account-group/{id} [put]
// @Security Bearer
func (e BusExchangeAccountGroup) Update(c *gin.Context) {
	req := dto.BusExchangeAccountGroupUpdateReq{}
	s := service.BusExchangeAccountGroup{}
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
		e.Error(500, err, fmt.Sprintf("修改账号组设置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除账号组设置
// @Summary 删除账号组设置
// @Description 删除账号组设置
// @Tags 账号组设置
// @Param data body dto.BusExchangeAccountGroupDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/account-group [delete]
// @Security Bearer
func (e BusExchangeAccountGroup) Delete(c *gin.Context) {
	s := service.BusExchangeAccountGroup{}
	req := dto.BusExchangeAccountGroupDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除账号组设置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

func (e BusExchangeAccountGroup) GetAccountGroupListByAccountId(c *gin.Context) {
	s := service.BusExchangeAccountGroup{}
	req := dto.BusAccountGroupListGetReq{}
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

	list := make([]models.BusExchangeAccountGroup, 0)
	err = s.GetGroupListByAccountId(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取账户绑定账户组失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "查询成功")
}
