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

type BusStrategyConfigDict struct {
	api.Api
}

// GetPage 获取策略配置字典列表
// @Summary 获取策略配置字典列表
// @Description 获取策略配置字典列表
// @Tags 策略配置字典
// @Param strategyId query string false "id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusStrategyConfigDict}} "{"code": 200, "data": [...]}"
// @Router /api/v1/strategy-config [get]
// @Security Bearer
func (e BusStrategyConfigDict) GetPage(c *gin.Context) {
	req := dto.BusStrategyConfigDictGetPageReq{}
	s := service.BusStrategyConfigDict{}
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
	list := make([]models.BusStrategyConfigDict, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取策略配置字典失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取策略配置字典
// @Summary 获取策略配置字典
// @Description 获取策略配置字典
// @Tags 策略配置字典
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusStrategyConfigDict} "{"code": 200, "data": [...]}"
// @Router /api/v1/strategy-config/{id} [get]
// @Security Bearer
func (e BusStrategyConfigDict) Get(c *gin.Context) {
	req := dto.BusStrategyConfigDictGetReq{}
	s := service.BusStrategyConfigDict{}
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
	var object models.BusStrategyConfigDict

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取策略配置字典失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建策略配置字典
// @Summary 创建策略配置字典
// @Description 创建策略配置字典
// @Tags 策略配置字典
// @Accept application/json
// @Product application/json
// @Param data body dto.BusStrategyConfigDictInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/strategy-config [post]
// @Security Bearer
func (e BusStrategyConfigDict) Insert(c *gin.Context) {
	req := dto.BusStrategyConfigDictInsertReq{}
	s := service.BusStrategyConfigDict{}
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
		e.Error(500, err, fmt.Sprintf("创建策略配置字典失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改策略配置字典
// @Summary 修改策略配置字典
// @Description 修改策略配置字典
// @Tags 策略配置字典
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusStrategyConfigDictUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/strategy-config/{id} [put]
// @Security Bearer
func (e BusStrategyConfigDict) Update(c *gin.Context) {
	req := dto.BusStrategyConfigDictUpdateReq{}
	s := service.BusStrategyConfigDict{}
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
		e.Error(500, err, fmt.Sprintf("修改策略配置字典失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除策略配置字典
// @Summary 删除策略配置字典
// @Description 删除策略配置字典
// @Tags 策略配置字典
// @Param data body dto.BusStrategyConfigDictDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/strategy-config [delete]
// @Security Bearer
func (e BusStrategyConfigDict) Delete(c *gin.Context) {
	s := service.BusStrategyConfigDict{}
	req := dto.BusStrategyConfigDictDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除策略配置字典失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
