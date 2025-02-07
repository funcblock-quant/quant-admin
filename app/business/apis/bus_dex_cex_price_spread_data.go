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

type BusDexCexPriceSpreadData struct {
	api.Api
}

// GetPage 获取dex-cex价差数据列表
// @Summary 获取dex-cex价差数据列表
// @Description 获取dex-cex价差数据列表
// @Tags dex-cex价差数据
// @Param symbol query string false "观察币种"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusDexCexPriceSpreadData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-cex-price-spread-data [get]
// @Security Bearer
func (e BusDexCexPriceSpreadData) GetPage(c *gin.Context) {
	req := dto.BusDexCexPriceSpreadDataGetPageReq{}
	s := service.BusDexCexPriceSpreadData{}
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
	list := make([]models.BusDexCexPriceSpreadData, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取dex-cex价差数据失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetDexCexHistoryChart 获取dex-cex价差历史数据-图表形式
// @Summary 获取dex-cex价差历史数据-图表形式
// @Description 获取dex-cex价差历史数据-图表形式
// @Tags dex-cex价差图表数据
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.BusDexCexTriangularSpreadHistory}} "{"code": 200, "data": [...]}"
// @Router /api/v1/getDexCexHistoryChart [get]
// @Security Bearer
func (e BusDexCexPriceSpreadData) GetDexCexHistoryChart(c *gin.Context) {
	req := dto.BusDexCexPriceSpreadDataHistoryChartReq{}
	s := service.BusDexCexPriceSpreadData{}
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

	chart := dto.BusDexCexTriangularSpreadHistory{}

	err = s.GetDexCexHistoryChart(&req, &chart)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取dex-cex价差图表数据失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(chart, "查询成功")
}

// Get 获取dex-cex价差数据
// @Summary 获取dex-cex价差数据
// @Description 获取dex-cex价差数据
// @Tags dex-cex价差数据
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusDexCexPriceSpreadData} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-cex-price-spread-data/{id} [get]
// @Security Bearer
func (e BusDexCexPriceSpreadData) Get(c *gin.Context) {
	req := dto.BusDexCexPriceSpreadDataGetReq{}
	s := service.BusDexCexPriceSpreadData{}
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
	var object models.BusDexCexPriceSpreadData

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取dex-cex价差数据失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建dex-cex价差数据
// @Summary 创建dex-cex价差数据
// @Description 创建dex-cex价差数据
// @Tags dex-cex价差数据
// @Accept application/json
// @Product application/json
// @Param data body dto.BusDexCexPriceSpreadDataInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-dex-cex-price-spread-data [post]
// @Security Bearer
func (e BusDexCexPriceSpreadData) Insert(c *gin.Context) {
	req := dto.BusDexCexPriceSpreadDataInsertReq{}
	s := service.BusDexCexPriceSpreadData{}
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
		e.Error(500, err, fmt.Sprintf("创建dex-cex价差数据失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改dex-cex价差数据
// @Summary 修改dex-cex价差数据
// @Description 修改dex-cex价差数据
// @Tags dex-cex价差数据
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusDexCexPriceSpreadDataUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-dex-cex-price-spread-data/{id} [put]
// @Security Bearer
func (e BusDexCexPriceSpreadData) Update(c *gin.Context) {
	req := dto.BusDexCexPriceSpreadDataUpdateReq{}
	s := service.BusDexCexPriceSpreadData{}
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
		e.Error(500, err, fmt.Sprintf("修改dex-cex价差数据失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除dex-cex价差数据
// @Summary 删除dex-cex价差数据
// @Description 删除dex-cex价差数据
// @Tags dex-cex价差数据
// @Param data body dto.BusDexCexPriceSpreadDataDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-dex-cex-price-spread-data [delete]
// @Security Bearer
func (e BusDexCexPriceSpreadData) Delete(c *gin.Context) {
	s := service.BusDexCexPriceSpreadData{}
	req := dto.BusDexCexPriceSpreadDataDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除dex-cex价差数据失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
