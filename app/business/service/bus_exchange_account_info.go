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

type BusExchangeAccountInfo struct {
	service.Service
}

// GetPage 获取BusExchangeAccountInfo列表
func (e *BusExchangeAccountInfo) GetPage(c *dto.BusExchangeAccountInfoGetPageReq, p *actions.DataPermission, list *[]models.BusExchangeAccountInfo, count *int64) error {
	var err error
	var data models.BusExchangeAccountInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusExchangeAccountInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusExchangeAccountInfo对象
func (e *BusExchangeAccountInfo) Get(d *dto.BusExchangeAccountInfoGetReq, p *actions.DataPermission, model *models.BusExchangeAccountInfo) error {
	var data models.BusExchangeAccountInfo

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
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

// Insert 创建BusExchangeAccountInfo对象
func (e *BusExchangeAccountInfo) Insert(c *dto.BusExchangeAccountInfoInsertReq) error {
	var data models.BusExchangeAccountInfo
	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("Service InsertBusExchangeAccountGroup error:%s \r\n", tx.Error)
		return tx.Error
	}

	// 插入账户信息
	c.Generate(&data)
	err := tx.Create(&data).Error
	if err != nil {
		tx.Rollback() // 插入主表失败，回滚事务
		e.Log.Errorf("BusExchangeAccountInfoService Insert error:%s \r\n", err)
		return err
	}

	// 保存账户绑定的账户组
	accountId := data.Id
	accountGroupRelation := make([]models.BusExchangeAccountGroupRelation, 0)
	for _, groupId := range c.AccountGroupIds {
		relation := models.BusExchangeAccountGroupRelation{AccountId: strconv.Itoa(accountId), GroupId: groupId}
		accountGroupRelation = append(accountGroupRelation, relation)
	}

	if err := tx.CreateInBatches(accountGroupRelation, len(accountGroupRelation)).Error; err != nil {
		tx.Rollback() // 插入关系表失败，回滚事务
		e.Log.Errorf("Error while inserting BusExchangeAccountGroupRelation: %v", err)
		return err
	}
	// 3. 提交事务
	err = tx.Commit().Error
	return err
}

// Update 修改BusExchangeAccountInfo对象
func (e *BusExchangeAccountInfo) Update(c *dto.BusExchangeAccountInfoUpdateReq, p *actions.DataPermission) error {
	var data = models.BusExchangeAccountInfo{}

	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("Service Update BusExchangeAccountInfo error:%s \r\n", tx.Error)
		return tx.Error
	}

	tx.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := tx.Save(&data)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service BusExchangeAccountInfo Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 2. 先删除绑定的账户组，保存账户绑定的账户组
	accountId := data.Id
	tx.Where("account_id = ? and deleted_at is null", accountId)
	tx.Debug().Delete(&models.BusExchangeAccountGroupRelation{}, "account_id=?", accountId)
	if err := tx.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("BusExchangeAccountGroupService Save error:%s \r\n", err)
		return err
	}

	accountGroupRelation := make([]models.BusExchangeAccountGroupRelation, 0)
	for _, groupId := range c.AccountGroupIds {
		relation := models.BusExchangeAccountGroupRelation{AccountId: strconv.Itoa(accountId), GroupId: groupId}
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

	return nil
}

// Remove 删除BusExchangeAccountInfo
func (e *BusExchangeAccountInfo) Remove(d *dto.BusExchangeAccountInfoDeleteReq, p *actions.DataPermission) error {
	var data models.BusExchangeAccountInfo
	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("Service InsertBusExchangeAccountGroup error:%s \r\n", tx.Error)
		return tx.Error
	}

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Service RemoveBusExchangeAccountInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}

	// 删除账户绑定的所有账户组
	if err := tx.Debug().Where("account_id = ?", d.GetId()).Delete(&models.BusExchangeAccountGroupRelation{}).Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Service RemoveBusExchangeAccountGroupRelation error: %s \r\n", err)
		return err
	}

	// 3. 提交事务
	if err := tx.Commit().Error; err != nil {
		e.Log.Errorf("Service RemoveBusExchangeAccountInfo Commit error: %s \r\n", err)
		return err
	}

	return nil
}

// GetAccountListByGroupId  Get 根据groupid获取已绑定的BusExchangeAccountInfo对象
func (e *BusExchangeAccountInfo) GetAccountListByGroupId(d *dto.BusGroupAccountInfoGetReq, p *actions.DataPermission, list *[]models.BusExchangeAccountInfo) error {
	var groupRelation models.BusExchangeAccountGroupRelation
	var accountInfo models.BusExchangeAccountInfo

	// 子查询：获取所有未删除的 groupRelation 记录的 accountId
	subQuery := e.Orm.Model(&groupRelation).
		Select("account_id").
		Where("group_id = ? AND deleted_at IS NULL", d.GroupId)

	// 主查询：通过子查询的 accountId 查找 BusExchangeAccountInfo 的记录
	err := e.Orm.Model(&accountInfo).
		Scopes(
			actions.Permission(accountInfo.TableName(), p), // 数据权限处理
		).
		Where("id IN (?)", subQuery).
		Find(list).Error

	e.Log.Infof("GetAccountListByGroupId %+v\r\n", list)

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
