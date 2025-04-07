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

// GetExchangeList 获取观察机器人所有的交易所列表
func (e BusDexCexTriangularObserver) GetExchangeList(c *gin.Context) {
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
	list := make([]dto.DexCexTriangularObserverExchangeListResp, 0)

	err = s.GetExchangeList(p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下三角套利观察币种列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}

// GetDexWalletList 获取观察机器人所有的dex钱包列表
func (e BusDexCexTriangularObserver) GetDexWalletList(c *gin.Context) {
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
	list := make([]models.BusDexWallet, 0)

	err = s.GetDexWalletList(p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下三角套利钱包列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}

// GetCexAccountList 获取观察机器人所有的Cex账户列表
func (e BusDexCexTriangularObserver) GetCexAccountList(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusGetCexAccountListReq{}
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
	list := make([]models.BusExchangeAccountInfo, 0)

	err = s.GetCexAccountList(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下三角套利交易所账户列表失败，\r\n失败信息 %s", err.Error()))
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
	var object dto.BusDexCexTriangularObserverDetailResp

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

func (e BusDexCexTriangularObserver) StartTrader(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusDexCexTriangularObserverStartTraderReq{}
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
	err = s.StartTrader(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("启动交易失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "开启交易成功")
}

func (e BusDexCexTriangularObserver) StopTrader(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusDexCexTriangularObserverStopTraderReq{}
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
	e.Logger.Infof("User: %d Stop Trader req: %+v", req.UpdateBy, req)
	err = s.StopTrader(&req, false)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("停止交易失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "停止交易成功")
}

func (e BusDexCexTriangularObserver) UpdateObserver(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusDexCexTriangularUpdateObserverParamsReq{}
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
	err = s.UpdateObserver(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("更新observer 参数失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "更新成功")
}

func (e BusDexCexTriangularObserver) UpdateTrader(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusDexCexTriangularUpdateTraderParamsReq{}
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
	err = s.UpdateTrader(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("更新trader 参数失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "更新成功")
}

// StopAllTrades 一键暂停所有交易
func (e BusDexCexTriangularObserver) StopAllTrades(c *gin.Context) {
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
	userId := user.GetUserId(c)
	failedList, err := s.StopAllTrader(userId)
	if err != nil {
		if len(failedList) == 0 {
			e.Error(500, err, fmt.Sprintf("更新trader 参数失败，\r\n失败信息 %s", err.Error()))
		} else {
			// 失败，且有部分 trader 停止失败
			e.Error(500, err, fmt.Sprintf("停止trader部分成果，失败列表: %v", failedList))
		}
		return
	}

	e.OK(nil, "更新成功")
}

func (e BusDexCexTriangularObserver) UpdateWaterLevel(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusDexCexTriangularUpdateWaterLevelParamsReq{}
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
	err = s.UpdateWaterLevel(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("更新trader 参数失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "更新成功")
}

func (e BusDexCexTriangularObserver) GetGlobalWaterLevelState(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusGetCexExchangeConfigListReq{}
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
	resp, err := s.GetGlobalWaterLevelState(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取全局水位调节 参数失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(resp, "获取全局水位调节参数成功")
}

func (e BusDexCexTriangularObserver) UpdateGlobalWaterLevel(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusDexCexTriangularUpdateGlobalWaterLevelConfigReq{}
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
	err = s.UpdateGlobalWaterLevelConfigV2(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("更新全局水位调节 参数失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "更新成功")
}

func (e BusDexCexTriangularObserver) GetGlobalRiskConfigState(c *gin.Context) {
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
	resp, err := s.GetGlobalRiskConfigState()
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取全局风控 参数失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(resp, "获取全局风控参数成功")
}

func (e BusDexCexTriangularObserver) UpdateGlobalRiskConfig(c *gin.Context) {
	s := service.BusDexCexTriangularObserver{}
	req := dto.BusDexCexTriangularUpdateGlobalRiskConfig{}
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
	err = s.UpdateGlobalRiskConfig(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("更新全局风控 参数失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "更新成功")
}

func (e BusDexCexTriangularObserver) GetBoundAccountList(c *gin.Context) {
	req := dto.BusGetBoundAccountReq{}
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
	resp := dto.BusGetBoundAccountResp{}

	err = s.GetBoundAccountList(&req, p, &resp)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下账号绑定列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(resp, "查询成功")
}

func (e BusDexCexTriangularObserver) GetCanBoundAccountList(c *gin.Context) {
	req := dto.BusGetBoundAccountReq{}
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
	resp := dto.BusGetBoundAccountResp{}

	err = s.GetCanBoundAccountList(&req, p, &resp)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下账号可绑定列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(resp, "查询成功")
}

func (e BusDexCexTriangularObserver) GetActiveAccountPairs(c *gin.Context) {
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
	resp := make([]dto.BusAccountPairInfo, 0)

	err = s.GetActiveAccountPairs(p, &resp)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取链上链下账号可绑定列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(resp, "查询成功")
}

func (e BusDexCexTriangularObserver) GetRealtimeInterestRate(c *gin.Context) {
	req := dto.BusGetInterestRateReq{}
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
	resp := dto.BusGetInterestRateResp{}

	err = s.GetRealtimeInterestRate(&req, p, &resp)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取实时汇率失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(resp, "查询成功")
}
