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

type BusDexCexTriangularArbitrageRecord struct {
	api.Api
}

// GetPage 获取链上链下三角套利记录列表
// @Summary 获取链上链下三角套利记录列表
// @Description 获取链上链下三角套利记录列表
// @Tags 链上链下三角套利记录
// @Param strategyId query string false "策略id"
// @Param arbitrageId query string false "套利记录id"
// @Param type query string false "套利类型"
// @Param status query string false "套利单状态"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusDexCexTriangularArbitrageRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/dex-cex-triangular-record [get]
// @Security Bearer
func (e BusDexCexTriangularArbitrageRecord) GetPage(c *gin.Context) {
	req := dto.BusDexCexTriangularArbitrageRecordGetPageReq{}
	s := service.BusDexCexTriangularArbitrageRecord{}
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
	list := make([]models.BusDexCexTriangularArbitrageRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下三角套利记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取链上链下三角套利记录
// @Summary 获取链上链下三角套利记录
// @Description 获取链上链下三角套利记录
// @Tags 链上链下三角套利记录
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusDexCexTriangularArbitrageRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/dex-cex-triangular-record/{id} [get]
// @Security Bearer
func (e BusDexCexTriangularArbitrageRecord) Get(c *gin.Context) {
	req := dto.BusDexCexTriangularArbitrageRecordGetReq{}
	s := service.BusDexCexTriangularArbitrageRecord{}
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
	var object models.BusDexCexTriangularArbitrageRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下三角套利记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建链上链下三角套利记录
// @Summary 创建链上链下三角套利记录
// @Description 创建链上链下三角套利记录
// @Tags 链上链下三角套利记录
// @Accept application/json
// @Product application/json
// @Param data body dto.BusDexCexTriangularArbitrageRecordInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/dex-cex-triangular-record [post]
// @Security Bearer
func (e BusDexCexTriangularArbitrageRecord) Insert(c *gin.Context) {
	req := dto.BusDexCexTriangularArbitrageRecordInsertReq{}
	s := service.BusDexCexTriangularArbitrageRecord{}
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
		e.Error(500, err, fmt.Sprintf("创建链上链下三角套利记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改链上链下三角套利记录
// @Summary 修改链上链下三角套利记录
// @Description 修改链上链下三角套利记录
// @Tags 链上链下三角套利记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusDexCexTriangularArbitrageRecordUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/dex-cex-triangular-record/{id} [put]
// @Security Bearer
func (e BusDexCexTriangularArbitrageRecord) Update(c *gin.Context) {
	req := dto.BusDexCexTriangularArbitrageRecordUpdateReq{}
	s := service.BusDexCexTriangularArbitrageRecord{}
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
		e.Error(500, err, fmt.Sprintf("修改链上链下三角套利记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除链上链下三角套利记录
// @Summary 删除链上链下三角套利记录
// @Description 删除链上链下三角套利记录
// @Tags 链上链下三角套利记录
// @Param data body dto.BusDexCexTriangularArbitrageRecordDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/dex-cex-triangular-record [delete]
// @Security Bearer
func (e BusDexCexTriangularArbitrageRecord) Delete(c *gin.Context) {
	s := service.BusDexCexTriangularArbitrageRecord{}
	req := dto.BusDexCexTriangularArbitrageRecordDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除链上链下三角套利记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// GetArbitrageOpportunity 获取套利机会列表
func (e BusDexCexTriangularArbitrageRecord) GetArbitrageOpportunity(c *gin.Context) {
	req := dto.BusArbitrageOpportunityGetReq{}
	s := service.BusDexCexTriangularArbitrageRecord{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	log.Infof("request data:%+v\r\n", req)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	list := make([]dto.BusArbitrageOpportunityGetResp, 0)

	err = s.QueryArbitrageOpportunityList(&req, p, &list)
	e.Logger.Infof("list data:%+v\r\n", list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取套利机会失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}
