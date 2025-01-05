package service

import (
	"errors"
	"strconv"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusExchangeAccountGroup struct {
	service.Service
}

// GetPage 获取BusExchangeAccountGroup列表
func (e *BusExchangeAccountGroup) GetPage(c *dto.BusExchangeAccountGroupGetPageReq, p *actions.DataPermission, list *[]models.BusExchangeAccountGroup, count *int64) error {
	var err error
	var data models.BusExchangeAccountGroup

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusExchangeAccountGroupService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusExchangeAccountGroup对象
func (e *BusExchangeAccountGroup) Get(d *dto.BusExchangeAccountGroupGetReq, p *actions.DataPermission, model *models.BusExchangeAccountGroup) error {
	var data models.BusExchangeAccountGroup

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusExchangeAccountGroup error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusExchangeAccountGroup对象
func (e *BusExchangeAccountGroup) Insert(c *dto.BusExchangeAccountGroupInsertReq) error {
	var data models.BusExchangeAccountGroup

	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("Service InsertBusExchangeAccountGroup error:%s \r\n", tx.Error)
		return tx.Error
	}
	// 1. 生成主表数据，并插入主表
	c.Generate(&data)
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback() // 插入主表失败，回滚事务
		e.Log.Errorf("Error while inserting InsertBusExchangeAccountGroup: %v", err)
		return err
	}

	// 2. 保存账户组下绑定的账户
	groupId := data.Id
	accountGroupRelation := make([]models.BusExchangeAccountGroupRelation, 0)
	for _, accountId := range c.AccountIds {
		relation := models.BusExchangeAccountGroupRelation{AccountId: strconv.Itoa(accountId), GroupId: strconv.Itoa(groupId)}
		accountGroupRelation = append(accountGroupRelation, relation)
	}
	if err := tx.CreateInBatches(accountGroupRelation, len(accountGroupRelation)).Error; err != nil {
		tx.Rollback() // 插入关系表失败，回滚事务
		e.Log.Errorf("Error while inserting BusExchangeAccountGroupRelation: %v", err)
		return err
	}
	// 3. 提交事务
	var err = tx.Commit().Error
	return err
}

// Update 修改BusExchangeAccountGroup对象
func (e *BusExchangeAccountGroup) Update(c *dto.BusExchangeAccountGroupUpdateReq, p *actions.DataPermission) error {
	var data = models.BusExchangeAccountGroup{}
	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("Service Update BusExchangeAccountGroup error:%s \r\n", tx.Error)
		return tx.Error
	}

	tx.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := tx.Save(&data)
	if err := db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Service BusExchangeAccountGroup Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 2. 先删除绑定的账户，保存账户组下绑定的账户
	groupId := data.Id
	tx.Where("group_id = ? and deleted_at is null", groupId)
	tx.Debug().Delete(&models.BusExchangeAccountGroupRelation{}, "group_id=?", groupId)
	if err := tx.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("BusExchangeAccountGroupService Save error:%s \r\n", err)
		return err
	}

	accountGroupRelation := make([]models.BusExchangeAccountGroupRelation, 0)
	for _, accountId := range c.AccountIds {
		relation := models.BusExchangeAccountGroupRelation{AccountId: strconv.Itoa(accountId), GroupId: strconv.Itoa(groupId)}
		accountGroupRelation = append(accountGroupRelation, relation)
	}
	if err := tx.CreateInBatches(accountGroupRelation, len(accountGroupRelation)).Error; err != nil {
		tx.Rollback() // 插入关系表失败，回滚事务
		e.Log.Errorf("Error while inserting BusExchangeAccountGroupRelation: %v", err)
		return err
	}
	// 3. 提交事务
	var err = tx.Commit().Error
	return err
}

// Remove 删除BusExchangeAccountGroup
func (e *BusExchangeAccountGroup) Remove(d *dto.BusExchangeAccountGroupDeleteReq, p *actions.DataPermission) error {
	var data models.BusExchangeAccountGroup
	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("Service InsertBusExchangeAccountGroup error:%s \r\n", tx.Error)
		return tx.Error
	}
	// 删除主表数据
	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Service RemoveBusExchangeAccountGroup error:%s \r\n", err)
		return err
	}

	// 如果没有权限删除任何数据
	if db.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("无权删除该数据")
	}

	// 删除账户组绑定的账户
	if err := tx.Debug().Where("group_id = ?", d.GetId()).Delete(&models.BusExchangeAccountGroupRelation{}).Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Service RemoveBusExchangeAccountGroupRelation error: %s \r\n", err)
		return err
	}

	// 3. 提交事务
	if err := tx.Commit().Error; err != nil {
		e.Log.Errorf("Service RemoveBusExchangeAccountGroup Commit error: %s \r\n", err)
		return err
	}

	return nil
}

// GetGroupListByAccountId  Get 根据account获取已绑定的BusExchangeAccountGroup对象
func (e *BusExchangeAccountGroup) GetGroupListByAccountId(d *dto.BusAccountGroupListGetReq, p *actions.DataPermission, list *[]models.BusExchangeAccountGroup) error {
	var groupRelation models.BusExchangeAccountGroupRelation
	var groupInfo models.BusExchangeAccountGroup

	// 子查询：获取所有未删除的 groupRelation 记录的 groupId
	subQuery := e.Orm.Model(&groupRelation).
		Select("group_id").
		Where("account_id = ? AND deleted_at IS NULL", d.AccountId)

	// 主查询：通过子查询的 groupId 查找 BusExchangeAccountGroup 的记录
	err := e.Orm.Model(&groupInfo).
		Scopes(
			actions.Permission(groupInfo.TableName(), p), // 数据权限处理
		).
		Where("id IN (?)", subQuery).
		Find(list).Error

	e.Log.Infof("GetGroupListByAccountId %+v\r\n", list)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusExchangeAccountInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
