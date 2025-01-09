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

type BusDexCexTriangularObserver struct {
	api.Api
}

// GetPage 获取链上链下三角套利观察列表
// @Summary 获取链上链下三角套利观察列表
// @Description 获取链上链下三角套利观察列表
// @Tags 链上链下三角套利观察
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusDexCexTriangularObserver}} "{"code": 200, "data": [...]}"
// @Router /api/v1/dex-cex-triangular-observer [get]
// @Security Bearer
func (e BusDexCexTriangularObserver) GetPage(c *gin.Context) {
	req := dto.BusDexCexTriangularObserverGetPageReq{}
	s := service.BusDexCexTriangularObserver{}
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
	list := make([]dto.BusDexCexTriangularObserverGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下三角套利观察失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetSymbolList 获取观察机器人所有的币种列表
func (e BusDexCexTriangularObserver) GetSymbolList(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]dto.DexCexTriangularObserverSymbolListResp, 0)

	err = s.GetSymbolList(p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下三角套利观察币种列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}

// Get 获取链上链下三角套利观察
// @Summary 获取链上链下三角套利观察
// @Description 获取链上链下三角套利观察
// @Tags 链上链下三角套利观察
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusDexCexTriangularObserver} "{"code": 200, "data": [...]}"
// @Router /api/v1/dex-cex-triangular-observer/{id} [get]
// @Security Bearer
func (e BusDexCexTriangularObserver) Get(c *gin.Context) {
	req := dto.BusDexCexTriangularObserverGetReq{}
	s := service.BusDexCexTriangularObserver{}
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
	var object models.BusDexCexTriangularObserver

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下三角套利观察失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建链上链下三角套利观察
// @Summary 创建链上链下三角套利观察
// @Description 创建链上链下三角套利观察
// @Tags 链上链下三角套利观察
// @Accept application/json
// @Product application/json
// @Param data body dto.BusDexCexTriangularObserverInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/dex-cex-triangular-observer [post]
// @Security Bearer
func (e BusDexCexTriangularObserver) Insert(c *gin.Context) {
	req := dto.BusDexCexTriangularObserverInsertReq{}
	s := service.BusDexCexTriangularObserver{}
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
		e.Error(500, err, fmt.Sprintf("创建链上链下三角套利观察失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// BatchInsert 批量创建链上链下三角套利观察
// @Summary 创建链上链下三角套利观察
// @Description 创建链上链下三角套利观察
// @Tags 链上链下三角套利观察
// @Accept application/json
// @Product application/json
// @Param data body dto.BusDexCexTriangularObserverInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/dex-cex-triangular-observer [post]
// @Security Bearer
func (e BusDexCexTriangularObserver) BatchInsert(c *gin.Context) {
	req := dto.BusDexCexTriangularObserverBatchInsertReq{}
	s := service.BusDexCexTriangularObserver{}
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

	err = s.BatchInsert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建链上链下三角套利观察失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "创建成功")
}

// Update 修改链上链下三角套利观察
// @Summary 修改链上链下三角套利观察
// @Description 修改链上链下三角套利观察
// @Tags 链上链下三角套利观察
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusDexCexTriangularObserverUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/dex-cex-triangular-observer/{id} [put]
// @Security Bearer
func (e BusDexCexTriangularObserver) Update(c *gin.Context) {
	req := dto.BusDexCexTriangularObserverUpdateReq{}
	s := service.BusDexCexTriangularObserver{}
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
		e.Error(500, err, fmt.Sprintf("修改链上链下三角套利观察失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除链上链下三角套利观察
// @Summary 删除链上链下三角套利观察
// @Description 删除链上链下三角套利观察
// @Tags 链上链下三角套利观察
// @Param data body dto.BusDexCexTriangularObserverDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/dex-cex-triangular-observer [delete]
// @Security Bearer
func (e BusDexCexTriangularObserver) Delete(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusDexCexTriangularObserverDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除链上链下三角套利观察失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
