package apis

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
)

type BusPriceTriggerStrategyApikeyConfig struct {
	api.Api
}

// GetPage 获取价格触发下单策略实例amber配置信息列表
// @Summary 获取价格触发下单策略实例amber配置信息列表
// @Description 获取价格触发下单策略实例amber配置信息列表
// @Tags 价格触发下单策略实例amber配置信息
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusPriceTriggerStrategyApikeyConfig}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-price-trigger-apikey-config [get]
// @Security Bearer
func (e BusPriceTriggerStrategyApikeyConfig) GetPage(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyApikeyConfigGetPageReq{}
	s := service.BusPriceTriggerStrategyApikeyConfig{}
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
	roleName := user.GetRoleName(c)
	var userId int
	if roleName != "admin" {
		e.Logger.Debugf("admin user id is: %d", userId)
		userId = user.GetUserId(c)
	}
	p := actions.GetPermissionFromContext(c)
	list := make([]models.BusPriceTriggerStrategyApikeyConfig, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count, &userId)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取价格触发下单策略实例amber配置信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取价格触发下单策略实例amber配置信息
// @Summary 获取价格触发下单策略实例amber配置信息
// @Description 获取价格触发下单策略实例amber配置信息
// @Tags 价格触发下单策略实例amber配置信息
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusPriceTriggerStrategyApikeyConfig} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-price-trigger-apikey-config/{id} [get]
// @Security Bearer
func (e BusPriceTriggerStrategyApikeyConfig) Get(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyApikeyConfigGetReq{}
	s := service.BusPriceTriggerStrategyApikeyConfig{}
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
	var object models.BusPriceTriggerStrategyApikeyConfig

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取价格触发下单策略实例amber配置信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建价格触发下单策略实例amber配置信息
// @Summary 创建价格触发下单策略实例amber配置信息
// @Description 创建价格触发下单策略实例amber配置信息
// @Tags 价格触发下单策略实例amber配置信息
// @Accept application/json
// @Product application/json
// @Param data body dto.BusPriceTriggerStrategyApikeyConfigInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-price-trigger-apikey-config [post]
// @Security Bearer
func (e BusPriceTriggerStrategyApikeyConfig) Insert(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyApikeyConfigInsertReq{}
	s := service.BusPriceTriggerStrategyApikeyConfig{}
	userId := user.GetUserId(c)
	req.UserId = strconv.Itoa(userId)
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
		e.Error(500, err, fmt.Sprintf("创建价格触发下单策略实例amber配置信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// CheckApiKeyHealth 检查API key正确性
func (e BusPriceTriggerStrategyApikeyConfig) CheckApiKeyHealth(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyApikeyConfigCheckReq{}
	s := service.BusPriceTriggerStrategyApikeyConfig{}
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
	isHealth, err := s.CheckApiKeyHealth(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("测试api key连接失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(isHealth, "检查完成")
}

// Update 修改价格触发下单策略实例amber配置信息
// @Summary 修改价格触发下单策略实例amber配置信息
// @Description 修改价格触发下单策略实例amber配置信息
// @Tags 价格触发下单策略实例amber配置信息
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusPriceTriggerStrategyApikeyConfigUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-price-trigger-apikey-config/{id} [put]
// @Security Bearer
func (e BusPriceTriggerStrategyApikeyConfig) Update(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyApikeyConfigUpdateReq{}
	s := service.BusPriceTriggerStrategyApikeyConfig{}
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
		e.Error(500, err, fmt.Sprintf("修改价格触发下单策略实例amber配置信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除价格触发下单策略实例amber配置信息
// @Summary 删除价格触发下单策略实例amber配置信息
// @Description 删除价格触发下单策略实例amber配置信息
// @Tags 价格触发下单策略实例amber配置信息
// @Param data body dto.BusPriceTriggerStrategyApikeyConfigDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-price-trigger-apikey-config [delete]
// @Security Bearer
func (e BusPriceTriggerStrategyApikeyConfig) Delete(c *gin.Context) {
	s := service.BusPriceTriggerStrategyApikeyConfig{}
	req := dto.BusPriceTriggerStrategyApikeyConfigDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除价格触发下单策略实例amber配置信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
