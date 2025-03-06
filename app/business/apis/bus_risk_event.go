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

type BusRiskEvent struct {
	api.Api
}

// GetPage 获取风控记录列表
// @Summary 获取风控记录列表
// @Description 获取风控记录列表
// @Tags 风控记录
// @Param strategyId query string false "策略id"
// @Param strategyInstanceId query string false "策略实例ID"
// @Param riskScope query string false "风控范围"
// @Param assetSymbol query string false "风控币种"
// @Param riskLevel query string false "风控级别"
// @Param manualRecover query string false "是否需要人工恢复"
// @Param autoRecoverTime query time.Time false "自动恢复时间"
// @Param isRecovered query string false "是否已恢复"
// @Param recoveredBy query string false "审核人"
// @Param recoveredAt query time.Time false "恢复时间"
// @Param triggerRule query string false "触发的风控规则"
// @Param triggerValue query string false "触发值"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusRiskEvent}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-risk-event [get]
// @Security Bearer
func (e BusRiskEvent) GetPage(c *gin.Context) {
	req := dto.BusRiskEventGetPageReq{}
	s := service.BusRiskEvent{}
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
	list := make([]models.BusRiskEvent, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取风控记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取风控记录
// @Summary 获取风控记录
// @Description 获取风控记录
// @Tags 风控记录
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusRiskEvent} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-risk-event/{id} [get]
// @Security Bearer
func (e BusRiskEvent) Get(c *gin.Context) {
	req := dto.BusRiskEventGetReq{}
	s := service.BusRiskEvent{}
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
	var object models.BusRiskEvent

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取风控记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建风控记录
// @Summary 创建风控记录
// @Description 创建风控记录
// @Tags 风控记录
// @Accept application/json
// @Product application/json
// @Param data body dto.BusRiskEventInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-risk-event [post]
// @Security Bearer
func (e BusRiskEvent) Insert(c *gin.Context) {
	req := dto.BusRiskEventInsertReq{}
	s := service.BusRiskEvent{}
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
		e.Error(500, err, fmt.Sprintf("创建风控记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 恢复风控记录
// @Summary 恢复风控记录
// @Description 恢复风控记录
// @Tags 风控记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusRiskEventUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-risk-event/{id} [put]
// @Security Bearer
func (e BusRiskEvent) Update(c *gin.Context) {
	req := dto.BusRiskEventUpdateReq{}
	s := service.BusRiskEvent{}
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
		e.Error(500, err, fmt.Sprintf("修改风控记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除风控记录
// @Summary 删除风控记录
// @Description 删除风控记录
// @Tags 风控记录
// @Param data body dto.BusRiskEventDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-risk-event [delete]
// @Security Bearer
func (e BusRiskEvent) Delete(c *gin.Context) {
	s := service.BusRiskEvent{}
	req := dto.BusRiskEventDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除风控记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
