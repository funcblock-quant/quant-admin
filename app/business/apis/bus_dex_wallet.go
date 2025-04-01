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

type BusDexWallet struct {
	api.Api
}

// GetPage 获取dex钱包列表
// @Summary 获取dex钱包列表
// @Description 获取dex钱包列表
// @Tags dex钱包
// @Param walletName query string false "钱包名称"
// @Param walletAddress query string false "钱包地址"
// @Param blockchain query string false "链网络"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusDexWallet}} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-wallet [get]
// @Security Bearer
func (e BusDexWallet) GetPage(c *gin.Context) {
    req := dto.BusDexWalletGetPageReq{}
    s := service.BusDexWallet{}
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
	list := make([]models.BusDexWallet, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取dex钱包失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取dex钱包
// @Summary 获取dex钱包
// @Description 获取dex钱包
// @Tags dex钱包
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusDexWallet} "{"code": 200, "data": [...]}"
// @Router /api/v1/bus-dex-wallet/{id} [get]
// @Security Bearer
func (e BusDexWallet) Get(c *gin.Context) {
	req := dto.BusDexWalletGetReq{}
	s := service.BusDexWallet{}
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
	var object models.BusDexWallet

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取dex钱包失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建dex钱包
// @Summary 创建dex钱包
// @Description 创建dex钱包
// @Tags dex钱包
// @Accept application/json
// @Product application/json
// @Param data body dto.BusDexWalletInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/bus-dex-wallet [post]
// @Security Bearer
func (e BusDexWallet) Insert(c *gin.Context) {
    req := dto.BusDexWalletInsertReq{}
    s := service.BusDexWallet{}
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
		e.Error(500, err, fmt.Sprintf("创建dex钱包失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改dex钱包
// @Summary 修改dex钱包
// @Description 修改dex钱包
// @Tags dex钱包
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.BusDexWalletUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/bus-dex-wallet/{id} [put]
// @Security Bearer
func (e BusDexWallet) Update(c *gin.Context) {
    req := dto.BusDexWalletUpdateReq{}
    s := service.BusDexWallet{}
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
		e.Error(500, err, fmt.Sprintf("修改dex钱包失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除dex钱包
// @Summary 删除dex钱包
// @Description 删除dex钱包
// @Tags dex钱包
// @Param data body dto.BusDexWalletDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/bus-dex-wallet [delete]
// @Security Bearer
func (e BusDexWallet) Delete(c *gin.Context) {
    s := service.BusDexWallet{}
    req := dto.BusDexWalletDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除dex钱包失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
