package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusServerInfo struct {
	service.Service
}

// GetPage 获取BusServerInfo列表
func (e *BusServerInfo) GetPage(c *dto.BusServerInfoGetPageReq, p *actions.DataPermission, list *[]models.BusServerInfo, count *int64) error {
	var err error
	var data models.BusServerInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusServerInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusServerInfo对象
func (e *BusServerInfo) Get(d *dto.BusServerInfoGetReq, p *actions.DataPermission, model *models.BusServerInfo) error {
	var data models.BusServerInfo

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusServerInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusServerInfo对象
func (e *BusServerInfo) Insert(c *dto.BusServerInfoInsertReq) error {
	var err error
	var data models.BusServerInfo
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusServerInfoService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusServerInfo对象
func (e *BusServerInfo) Update(c *dto.BusServerInfoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusServerInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusServerInfoService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// StartServer 修改BusServerInfo对象
func (e *BusServerInfo) StartServer(c *dto.BusServerStartReq, p *actions.DataPermission) error {
	var err error
	var data = &models.BusServerInfo{}
	db := e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).Exec("UPDATE bus_server_info SET status = ? WHERE id = ?", 1, c.Id)

	if err = db.Error; err != nil {
		e.Log.Errorf("BusServerInfoService Stop error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// StopServer 修改BusServerInfo对象
func (e *BusServerInfo) StopServer(c *dto.BusServerStopReq, p *actions.DataPermission) error {
	var err error
	var data = &models.BusServerInfo{}
	db := e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).Exec("UPDATE bus_server_info SET status = ? WHERE id = ?", 2, c.Id)

	if err = db.Error; err != nil {
		e.Log.Errorf("BusServerInfoService Stop error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusServerInfo
func (e *BusServerInfo) Remove(d *dto.BusServerInfoDeleteReq, p *actions.DataPermission) error {
	var data models.BusServerInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusServerInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
