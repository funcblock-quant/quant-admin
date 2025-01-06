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

type BusStrategySymbolGroup struct {
	api.Api
}

// GetPage 获取观察组配置列表
// @Summary 获取观察组配置列表
// @Description 获取观察组配置列表
// @Tags 观察组配置
// @Param strategyInstanceId query string false "策略实例id"
// @Param groupType query string false "组类型"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusStrategySymbolGroup}} "{"code": 200, "data": [...]}"
// @Router /api/v1/symbol-group [get]
// @Security Bearer
func (e BusStrategySymbolGroup) GetPage(c *gin.Context) {
	req := dto.BusStrategySymbolGroupGetPageReq{}
	s := service.BusStrategySymbolGroup{}
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
	list := make([]models.BusStrategySymbolGroup, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取观察组配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取观察组配置
// @Summary 获取观察组配置
// @Description 获取观察组配置
// @Tags 观察组配置
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusStrategySymbolGroup} "{"code": 200, "data": [...]}"
// @Router /api/v1/symbol-group/{id} [get]
// @Security Bearer
func (e BusStrategySymbolGroup) Get(c *gin.Context) {
	req := dto.BusStrategySymbolGroupGetReq{}
	s := service.BusStrategySymbolGroup{}
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
	var object models.BusStrategySymbolGroup

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取观察组配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建观察组配置
// @Summary 创建观察组配置
// @Description 创建观察组配置
// @Tags 观察组配置
// @Accept application/json
// @Product application/json
// @Param data body dto.BusStrategySymbolGroupInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/symbol-group [post]
// @Security Bearer
func (e BusStrategySymbolGroup) Insert(c *gin.Context) {
	req := dto.BusStrategySymbolGroupInsertReq{}
	s := service.BusStrategySymbolGroup{}
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
		e.Error(500, err, fmt.Sprintf("创建观察组配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改观察组配置
// @Summary 修改观察组配置
// @Description 修改观察组配置
// @Tags 观察组配置
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusStrategySymbolGroupUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/symbol-group/{id} [put]
// @Security Bearer
func (e BusStrategySymbolGroup) Update(c *gin.Context) {
	req := dto.BusStrategySymbolGroupUpdateReq{}
	s := service.BusStrategySymbolGroup{}
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
		e.Error(500, err, fmt.Sprintf("修改观察组配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除观察组配置
// @Summary 删除观察组配置
// @Description 删除观察组配置
// @Tags 观察组配置
// @Param data body dto.BusStrategySymbolGroupDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/symbol-group [delete]
// @Security Bearer
func (e BusStrategySymbolGroup) Delete(c *gin.Context) {
	s := service.BusStrategySymbolGroup{}
	req := dto.BusStrategySymbolGroupDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除观察组配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
