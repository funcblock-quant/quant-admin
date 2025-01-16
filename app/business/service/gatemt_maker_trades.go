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

type GatemtMakerTrades struct {
	service.Service
}

// GetPage 获取GatemtMakerTrades列表
func (e *GatemtMakerTrades) GetPage(c *dto.GatemtMakerTradesGetPageReq, p *actions.DataPermission, list *[]models.GatemtMakerTrades, count *int64) error {
	var err error
	var data models.GatemtMakerTrades

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GatemtMakerTradesService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GatemtMakerTrades对象
func (e *GatemtMakerTrades) Get(d *dto.GatemtMakerTradesGetReq, p *actions.DataPermission, model *models.GatemtMakerTrades) error {
	var data models.GatemtMakerTrades

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGatemtMakerTrades error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GatemtMakerTrades对象
func (e *GatemtMakerTrades) Insert(c *dto.GatemtMakerTradesInsertReq) error {
	var err error
	var data models.GatemtMakerTrades
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GatemtMakerTradesService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GatemtMakerTrades对象
func (e *GatemtMakerTrades) Update(c *dto.GatemtMakerTradesUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.GatemtMakerTrades{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("GatemtMakerTradesService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除GatemtMakerTrades
func (e *GatemtMakerTrades) Remove(d *dto.GatemtMakerTradesDeleteReq, p *actions.DataPermission) error {
	var data models.GatemtMakerTrades

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveGatemtMakerTrades error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *GatemtMakerTrades) QueryTradesWithClientOrderId(d *dto.GatemtMakerTradesGetListReq, p *actions.DataPermission, list *[]models.GatemtMakerTrades) error {
	var err error
	var data models.GatemtMakerTrades

	err = e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("client_order_id = ?", d.ClientOrderId).
		Find(list).Limit(-1).Offset(-1).Error
	if err != nil {
		e.Log.Errorf("GatemtMakerTradesService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}
