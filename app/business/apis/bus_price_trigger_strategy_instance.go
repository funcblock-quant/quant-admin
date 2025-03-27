package apis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
)

type BusPriceTriggerStrategyInstance struct {
	api.Api
}

// GetPage 获取价格触发下单策略列表
// @Summary 获取价格触发下单策略列表
// @Description 获取价格触发下单策略列表
// @Tags 价格触发下单策略
// @Param closeTime query time.Time false "停止时间"
// @Param status query string false "状态，created, started, stopped, closed"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusPriceTriggerStrategyInstance}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-price-trigger-strategy [get]
// @Security Bearer
func (e BusPriceTriggerStrategyInstance) GetPage(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyInstanceGetPageReq{}
	s := service.BusPriceTriggerStrategyInstance{}
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
	userId := user.GetUserId(c)
	roleName := user.GetRoleName(c)
	roleNameList := strings.Split(roleName, ",")
	isAdminOrSystemAdmin := false
	for _, roleName := range roleNameList {
		if roleName == "admin" || roleName == "系统管理员" {
			isAdminOrSystemAdmin = true
			break
		}
	}

	if !isAdminOrSystemAdmin {
		//如果不是管理员，只能自己看自己添加的下单规则
		req.UserId = strconv.Itoa(userId)
	}
	p := actions.GetPermissionFromContext(c)
	list := make([]dto.BusPriceTriggerStrategyResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取价格触发下单策略失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取价格触发下单策略
// @Summary 获取价格触发下单策略
// @Description 获取价格触发下单策略
// @Tags 价格触发下单策略
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusPriceTriggerStrategyInstance} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-price-trigger-strategy/{id} [get]
// @Security Bearer
func (e BusPriceTriggerStrategyInstance) Get(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyInstanceGetReq{}
	s := service.BusPriceTriggerStrategyInstance{}
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
	var object models.BusPriceTriggerStrategyInstance

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取价格触发下单策略失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建价格触发下单策略
// @Summary 创建价格触发下单策略
// @Description 创建价格触发下单策略
// @Tags 价格触发下单策略
// @Accept application/json
// @Product application/json
// @Param data body dto.BusPriceTriggerStrategyInstanceInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-price-trigger-strategy [post]
// @Security Bearer
func (e BusPriceTriggerStrategyInstance) Insert(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyInstanceInsertReq{}
	s := service.BusPriceTriggerStrategyInstance{}
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
		e.Error(500, err, fmt.Sprintf("创建价格触发下单策略失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改价格触发下单策略
// @Summary 修改价格触发下单策略
// @Description 修改价格触发下单策略
// @Tags 价格触发下单策略
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusPriceTriggerStrategyInstanceUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-price-trigger-strategy/{id} [put]
// @Security Bearer
func (e BusPriceTriggerStrategyInstance) Update(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyInstanceUpdateReq{}
	s := service.BusPriceTriggerStrategyInstance{}
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
		e.Error(500, err, fmt.Sprintf("修改价格触发下单策略失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除价格触发下单策略
// @Summary 删除价格触发下单策略
// @Description 删除价格触发下单策略
// @Tags 价格触发下单策略
// @Param data body dto.BusPriceTriggerStrategyInstanceDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-price-trigger-strategy [delete]
// @Security Bearer
func (e BusPriceTriggerStrategyInstance) Delete(c *gin.Context) {
	s := service.BusPriceTriggerStrategyInstance{}
	req := dto.BusPriceTriggerStrategyInstanceDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除价格触发下单策略失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// StopInstance 暂停策略
// @Summary 暂停策略
// @Description 暂停策略
// @Tags 暂停策略
// @Accept application/json
// @Product application/json
// @Param data body dto.BusPriceTriggerStrategyInstanceInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "暂停实例成功"}"
// @Router /api/v1/stopTriggerInstance [post]
// @Security Bearer
func (e BusPriceTriggerStrategyInstance) StopInstance(c *gin.Context) {
	req := dto.StopTriggerInstanceRequest{}
	s := service.BusPriceTriggerStrategyInstance{}
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
	e.Logger.Infof("req:%#v", req)
	err = s.StopInstance(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("暂停实例失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK("", "暂停实例成功")
}

// UpdateProfitTarget 修改止盈配置
// @Summary 修改止盈配置
// @Description 修改止盈配置
// @Tags 修改止盈配置
// @Accept application/json
// @Product application/json
// @Param data body dto.BusPriceTriggerStrategyInstanceUpdateProfitTargetReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改止盈配置成功"}"
// @Router /api/v1/updateProfitTarget [put]
// @Security Bearer
func (e BusPriceTriggerStrategyInstance) UpdateProfitTarget(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyInstanceUpdateProfitTargetReq{}
	s := service.BusPriceTriggerStrategyInstance{}
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
	e.Logger.Infof("req:%#v", req)
	err = s.UpdateProfitTarget(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改止盈配置失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK("", "修改止盈配置成功")
}

// UpdateExecuteConfig 修改执行参数
// @Summary 修改执行参数
// @Description 修改执行参数
// @Tags 修改执行参数
// @Accept application/json
// @Product application/json
// @Param data body dto.BusPriceTriggerStrategyInstanceUpdateExecuteNumReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/updateExecuteConfig [put]
// @Security Bearer
func (e BusPriceTriggerStrategyInstance) UpdateExecuteConfig(c *gin.Context) {
	req := dto.BusPriceTriggerStrategyInstanceUpdateExecuteNumReq{}
	s := service.BusPriceTriggerStrategyInstance{}
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
	e.Logger.Infof("req:%#v", req)
	err = s.UpdateExecuteConfig(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改执行参数失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK("", "修改成功")
}

// GetSymbolList 获取价格触发下单策略实例所有的币种列表
func (e BusPriceTriggerStrategyInstance) GetSymbolList(c *gin.Context) {
	s := service.BusPriceTriggerStrategyInstance{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	list := make([]dto.BusPriceTriggerStrategySymbolListResp, 0)

	err = s.GetSymbolList(&list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取价格触发下单策略币种列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}

// GetExchangeUserIdList 获取价格触发下单策略实例所有的交易所userId列表
func (e BusPriceTriggerStrategyInstance) GetExchangeUserIdList(c *gin.Context) {
	s := service.BusPriceTriggerStrategyInstance{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	list := make([]dto.BusPriceTriggerStrategyExchangeUserIdListResp, 0)
	userId := user.GetUserId(c)
	err = s.GetExchangeUserIdList(userId, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取价格触发下单策略币种列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}
