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

type BusServerInfo struct {
	api.Api
}

// GetPage 获取服务器管理列表
// @Summary 获取服务器管理列表
// @Description 获取服务器管理列表
// @Tags 服务器管理
// @Param serverIp query string false "ip"
// @Param networkStatus query string false "网络健康"
// @Param status query string false "服务器状态"
// @Param cpuNum query string false "cpu核心数"
// @Param mermorySize query string false "内存大小"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusServerInfo}} "{"code": 200, "data": [...]}"
// @Router /api/v1/server-info [get]
// @Security Bearer
func (e BusServerInfo) GetPage(c *gin.Context) {
	req := dto.BusServerInfoGetPageReq{}
	s := service.BusServerInfo{}
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
	list := make([]models.BusServerInfo, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取服务器管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取服务器管理
// @Summary 获取服务器管理
// @Description 获取服务器管理
// @Tags 服务器管理
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusServerInfo} "{"code": 200, "data": [...]}"
// @Router /api/v1/server-info/{id} [get]
// @Security Bearer
func (e BusServerInfo) Get(c *gin.Context) {
	req := dto.BusServerInfoGetReq{}
	s := service.BusServerInfo{}
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
	var object models.BusServerInfo

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取服务器管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建服务器管理
// @Summary 创建服务器管理
// @Description 创建服务器管理
// @Tags 服务器管理
// @Accept application/json
// @Product application/json
// @Param data body dto.BusServerInfoInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/server-info [post]
// @Security Bearer
func (e BusServerInfo) Insert(c *gin.Context) {
	req := dto.BusServerInfoInsertReq{}
	s := service.BusServerInfo{}
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
		e.Error(500, err, fmt.Sprintf("创建服务器管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改服务器管理
// @Summary 修改服务器管理
// @Description 修改服务器管理
// @Tags 服务器管理
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusServerInfoUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/server-info/{id} [put]
// @Security Bearer
func (e BusServerInfo) Update(c *gin.Context) {
	req := dto.BusServerInfoUpdateReq{}
	s := service.BusServerInfo{}
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
		e.Error(500, err, fmt.Sprintf("修改服务器管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// StopServer 暂停服务器
// @Summary 暂停服务器
// @Description 暂停服务器
// @Tags 服务器管理
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusServerStopReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "暂停成功"}"
// @Router /api/v1/server-info/{id} [get]
// @Security Bearer
func (e BusServerInfo) StopServer(c *gin.Context) {
	req := dto.BusServerStopReq{}
	s := service.BusServerInfo{}
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

	err = s.StopServer(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("暂停服务器管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "暂停成功")
}

// StartServer 启用服务器
// @Summary 启用服务器
// @Description 启用服务器
// @Tags 服务器管理
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusServerStartReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "启用成功"}"
// @Router /api/v1/server-info/{id} [get]
// @Security Bearer
func (e BusServerInfo) StartServer(c *gin.Context) {
	req := dto.BusServerStartReq{}
	s := service.BusServerInfo{}
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

	err = s.StartServer(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("启用服务器管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "启用成功")
}

// Delete 删除服务器管理
// @Summary 删除服务器管理
// @Description 删除服务器管理
// @Tags 服务器管理
// @Param data body dto.BusServerInfoDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/server-info [delete]
// @Security Bearer
func (e BusServerInfo) Delete(c *gin.Context) {
	s := service.BusServerInfo{}
	req := dto.BusServerInfoDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除服务器管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
