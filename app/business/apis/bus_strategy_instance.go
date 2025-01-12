package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"quanta-admin/app/business/service"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
)

type BusStrategyInstance struct {
	api.Api
}

// GetPage 获取策略实例配置列表
// @Summary 获取策略实例配置列表
// @Description 获取策略实例配置列表
// @Tags 策略实例配置
// @Param strategyId query string false "策略id"
// @Param accountGroupId query string false "账户组id"
// @Param startRunTime query time.Time false "启动时间"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusStrategyInstance}} "{"code": 200, "data": [...]}"
// @Router /api/v1/strategy-instance [get]
// @Security Bearer
func (e BusStrategyInstance) GetPage(c *gin.Context) {
	req := dto.BusStrategyInstanceGetPageReq{}
	s := service.BusStrategyInstance{}
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
	list := make([]dto.BusStrategyInstanceGetPageResp, 0)
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
// @Success 200 {object} response.Response{data=models.BusStrategyInstance} "{"code": 200, "data": [...]}"
// @Router /api/v1/strategy-instance/{id} [get]
// @Security Bearer
func (e BusStrategyInstance) Get(c *gin.Context) {
	req := dto.BusStrategyInstanceGetReq{}
	s := service.BusStrategyInstance{}
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
	var object dto.BusStrategyInstanceGetResp

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
// @Param data body dto.BusStrategyInstanceInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/strategy-instance [post]
// @Security Bearer
func (e BusStrategyInstance) Insert(c *gin.Context) {
	req := dto.BusStrategyInstanceInsertReq{}
	s := service.BusStrategyInstance{}
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
// @Param data body dto.BusStrategyInstanceUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/strategy-instance/{id} [put]
// @Security Bearer
func (e BusStrategyInstance) Update(c *gin.Context) {
	req := dto.BusStrategyInstanceUpdateReq{}
	s := service.BusStrategyInstance{}
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

// StartInstance 启动策略实例
// @Summary 启动策略实例
// @Description 启动策略实例
// @Tags 策略实例启动
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusStrategyInstanceStartReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "启动成功"}"
// @Router /api/v1/startStrategyInstance/{id} [put]
// @Security Bearer
func (e BusStrategyInstance) StartInstance(c *gin.Context) {
	req := dto.BusStrategyInstanceStartReq{}
	s := service.BusStrategyInstance{}
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

	err = s.StartInstance(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("启动策略实例失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "启动成功")
}

// StopInstance 暂停策略实例
// @Summary 暂停策略实例
// @Description 暂停策略实例
// @Tags 策略实例暂停
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusStrategyInstanceStopReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "暂停成功"}"
// @Router /api/v1/stopStrategyInstance/{id} [put]
// @Security Bearer
func (e BusStrategyInstance) StopInstance(c *gin.Context) {
	req := dto.BusStrategyInstanceStopReq{}
	s := service.BusStrategyInstance{}
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

	err = s.StopInstance(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("暂停策略实例失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "暂停成功")
}

// Delete 删除策略实例配置
// @Summary 删除策略实例配置
// @Description 删除策略实例配置
// @Tags 策略实例配置
// @Param data body dto.BusStrategyInstanceDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/strategy-instance [delete]
// @Security Bearer
func (e BusStrategyInstance) Delete(c *gin.Context) {
	s := service.BusStrategyInstance{}
	req := dto.BusStrategyInstanceDeleteReq{}
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

// QueryInstanceDashboard 查询策略实例dashboard数据
// @Accept application/json
// @Product application/json
// @Param data body dto.BusStrategyInstanceInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/strategy-instance [post]
// @Security Bearer
func (e BusStrategyInstance) QueryInstanceDashboard(c *gin.Context) {
	req := dto.BusStrategyInstanceDashboardGetReq{}
	s := service.BusStrategyInstance{}
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
	resp := dto.BusStrategyInstanceDashboardGetResp{}
	err = s.QueryInstanceDashboard(&req, p, &resp)

	if err != nil {
		e.Error(500, err, fmt.Sprintf("查询策略dashboard失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(resp, "查询成功")
}
