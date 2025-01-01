package apis

import (
	"fmt"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service"
	"quanta-admin/app/business/service/dto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"quanta-admin/common/actions"
)

type BusArbitrageRecord struct {
	api.Api
}

// GetPage 获取套利记录表列表
// @Summary 获取套利记录表列表
// @Description 获取套利记录表列表
// @Tags 套利记录表
// @Param arbitrageId query string false "套利记录id"
// @Param strategyName query string false "策略名称"
// @Param contractType query string false "合约类型"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.BusArbitrageRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/arbitrage [get]
// @Security Bearer
func (e BusArbitrageRecord) GetPage(c *gin.Context) {
	req := dto.BusArbitrageRecordGetPageReq{}
	s := service.BusArbitrageRecord{}
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
	if req.BeginTime != "" && req.EndTime != "" {
		req.BeginTime, err = parseTimeToNs(e, req.BeginTime)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("时间格式错误，\r\n失败信息 %s", err.Error()))
		}
		req.EndTime, err = parseTimeToNs(e, req.EndTime)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("时间格式错误，\r\n失败信息 %s", err.Error()))
		}

	}
	p := actions.GetPermissionFromContext(c)
	list := make([]models.BusArbitrageRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取套利记录表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	// 转化为响应数据格式
	responseList := []BusArbitrageRecordResponse{}
	for _, record := range list {
		// 将查询结果转换为响应数据格式，包含格式化时间
		responseList = append(responseList, ToResponse(record))
	}

	e.PageOK(responseList, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func parseTimeToNs(e BusArbitrageRecord, timeStr string) (string, error) {
	// 定义时间格式，确保与前端传来的时间格式一致
	layout := "2006-01-02 15:04:05" // Go中的时间格式：2006年1月2日 15:04:05

	// 使用 time.Parse 函数将时间字符串解析为 time.Time
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return "", fmt.Errorf("时间格式不匹配：%w", err) // 使用 %w 将错误包裹后返回
	}

	// 转换为纳秒时间戳
	nsTimestamp := parsedTime.UnixNano()

	return strconv.FormatInt(nsTimestamp, 10), nil // 返回字符串形式的纳秒时间戳
}

// 格式化时间戳为字符串
func formatTimestampToString(ns int64) string {
	sec := ns / int64(time.Second)
	t := time.Unix(sec, 0)
	return t.Format("2006-01-02 15:04:05") // 可以根据需求更改时间格式
}

// 用于响应的结构体，添加时间格式字符串
type BusArbitrageRecordResponse struct {
	ArbitrageId        string `json:"arbitrageId"`
	StrategyInstanceId string `json:"strategyInstanceId"`
	StrategyName       string `json:"strategyName"`
	Type               string `json:"type"`
	ContractType       string `json:"contractType"`
	RealizedPnl        string `json:"realizedPnl"`
	UnrealizedPnl      string `json:"unrealizedPnl"`
	ExpectPnl          string `json:"expectPnl"`
	ExpectPnlPercent   string `json:"expectPnlPercent"`
	Status             string `json:"status"`
	StartTime          int64  `json:"startTime"`
	EndTime            int64  `json:"endTime"`
	// 时间字段转为格式字符串
	StartTimeStr string `json:"startTimeStr"`
	EndTimeStr   string `json:"endTimeStr"`
}

// 返回时，进行格式化处理，将时间戳转化为字符串
func ToResponse(e models.BusArbitrageRecord) BusArbitrageRecordResponse {
	return BusArbitrageRecordResponse{
		ArbitrageId:        e.ArbitrageId,
		StrategyInstanceId: e.StrategyInstanceId,
		StrategyName:       e.StrategyName,
		Type:               e.Type,
		ContractType:       e.ContractType,
		RealizedPnl:        e.RealizedPnl,
		UnrealizedPnl:      e.UnrealizedPnl,
		ExpectPnl:          e.ExpectPnl,
		ExpectPnlPercent:   e.ExpectPnlPercent,
		Status:             e.Status,
		StartTime:          e.StartTime,
		EndTime:            e.EndTime,
		// 格式化时间戳为字符串
		StartTimeStr: formatTimestampToString(e.StartTime),
		EndTimeStr:   formatTimestampToString(e.EndTime),
	}
}

// Get 获取套利记录表
// @Summary 获取套利记录表
// @Description 获取套利记录表
// @Tags 套利记录表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.BusArbitrageRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/arbitrage/{id} [get]
// @Security Bearer
func (e BusArbitrageRecord) Get(c *gin.Context) {
	req := dto.BusArbitrageRecordGetReq{}
	fmt.Printf("req1: %+v\n", req)
	s := service.BusArbitrageRecord{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	fmt.Printf("req2: %+v\n", req)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.BusArbitrageRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取套利记录表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

//// Insert 创建套利记录表
//// @Summary 创建套利记录表
//// @Description 创建套利记录表
//// @Tags 套利记录表
//// @Accept application/json
//// @Product application/json
//// @Param data body dto.BusArbitrageRecordInsertReq true "data"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
//// @Router /api/v1/arbitrage [post]
//// @Security Bearer
//func (e BusArbitrageRecord) Insert(c *gin.Context) {
//	req := dto.BusArbitrageRecordInsertReq{}
//	s := service.BusArbitrageRecord{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	// 设置创建人
//	req.SetCreateBy(user.GetUserId(c))
//
//	err = s.Insert(&req)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("创建套利记录表失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//
//	e.OK(req.GetId(), "创建成功")
//}
//
//// Update 修改套利记录表
//// @Summary 修改套利记录表
//// @Description 修改套利记录表
//// @Tags 套利记录表
//// @Accept application/json
//// @Product application/json
//// @Param id path int true "id"
//// @Param data body dto.BusArbitrageRecordUpdateReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
//// @Router /api/v1/arbitrage/{id} [put]
//// @Security Bearer
//func (e BusArbitrageRecord) Update(c *gin.Context) {
//	req := dto.BusArbitrageRecordUpdateReq{}
//	s := service.BusArbitrageRecord{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//	req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Update(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("修改套利记录表失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "修改成功")
//}
//
//// Delete 删除套利记录表
//// @Summary 删除套利记录表
//// @Description 删除套利记录表
//// @Tags 套利记录表
//// @Param data body dto.BusArbitrageRecordDeleteReq true "body"
//// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
//// @Router /api/v1/arbitrage [delete]
//// @Security Bearer
//func (e BusArbitrageRecord) Delete(c *gin.Context) {
//	s := service.BusArbitrageRecord{}
//	req := dto.BusArbitrageRecordDeleteReq{}
//	err := e.MakeContext(c).
//		MakeOrm().
//		Bind(&req).
//		MakeService(&s.Service).
//		Errors
//	if err != nil {
//		e.Logger.Error(err)
//		e.Error(500, err, err.Error())
//		return
//	}
//
//	// req.SetUpdateBy(user.GetUserId(c))
//	p := actions.GetPermissionFromContext(c)
//
//	err = s.Remove(&req, p)
//	if err != nil {
//		e.Error(500, err, fmt.Sprintf("删除套利记录表失败，\r\n失败信息 %s", err.Error()))
//		return
//	}
//	e.OK(req.GetId(), "删除成功")
//}
