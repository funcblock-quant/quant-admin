package apis

import (
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
)

type BusStrategyInstanceConfig struct {
	api.Api
}

// GetPage 获取策略实例配置列表
// @Summary 获取策略实例配置列表
// @Description 获取策略实例配置列表
// @Tags 策略实例配置
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusStrategyInstanceConfig}} "{"code": 200, "data": [...]}"
// @Router /api/v1/instance-config [get]
// @Security Bearer
func (e BusStrategyInstanceConfig) GetPage(c *gin.Context) {
	req := dto.BusStrategyInstanceConfigGetPageReq{}
	s := service.BusStrategyInstanceConfig{}
	log.Infof("api GetPage, req: %v", req)
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
	list := make([]models.BusStrategyInstanceConfig, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取策略实例配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取策略实例配置
// @Summary 获取策略实例配置
// @Description 获取策略实例配置
// @Tags 策略实例配置
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusStrategyInstanceConfig} "{"code": 200, "data": [...]}"
// @Router /api/v1/instance-config/{id} [get]
// @Security Bearer
func (e BusStrategyInstanceConfig) Get(c *gin.Context) {
	req := dto.BusStrategyInstanceConfigGetReq{}
	s := service.BusStrategyInstanceConfig{}
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
	var object models.BusStrategyInstanceConfig

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取策略实例配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建策略实例配置
// @Summary 创建策略实例配置
// @Description 创建策略实例配置
// @Tags 策略实例配置
// @Accept application/json
// @Product application/json
// @Param data body dto.BusStrategyInstanceConfigInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/instance-config [post]
// @Security Bearer
func (e BusStrategyInstanceConfig) Insert(c *gin.Context) {
	req := dto.BusStrategyInstanceConfigInsertReq{}
	s := service.BusStrategyInstanceConfig{}
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
		e.Error(500, err, fmt.Sprintf("创建策略实例配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改策略实例配置
// @Summary 修改策略实例配置
// @Description 修改策略实例配置
// @Tags 策略实例配置
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusStrategyInstanceConfigUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/instance-config/{id} [put]
// @Security Bearer
func (e BusStrategyInstanceConfig) Update(c *gin.Context) {
	req := dto.BusStrategyInstanceConfigUpdateReq{}
	s := service.BusStrategyInstanceConfig{}
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
		e.Error(500, err, fmt.Sprintf("修改策略实例配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除策略实例配置
// @Summary 删除策略实例配置
// @Description 删除策略实例配置
// @Tags 策略实例配置
// @Param data body dto.BusStrategyInstanceConfigDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/instance-config [delete]
// @Security Bearer
func (e BusStrategyInstanceConfig) Delete(c *gin.Context) {
	s := service.BusStrategyInstanceConfig{}
	req := dto.BusStrategyInstanceConfigDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除策略实例配置失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
