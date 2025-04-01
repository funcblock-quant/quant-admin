package service

import (
	"errors"
	"strconv"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/app/grpc/client"
	waterLevelPb "quanta-admin/app/grpc/proto/client/water_level_service"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
	"quanta-admin/common/global"
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
	e.Log.Debugf("Service Update BusExchangeAccountInfo data:%+v \r\n", c)
	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("Service Update BusExchangeAccountInfo error:%s \r\n", tx.Error)
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

	// 3. 提交事务
	var err = tx.Commit().Error
	return err

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

// QueryExchangeListInUse  获取BusExchangeAccountInfo所有关联的ExchangeList
func (e *BusExchangeAccountInfo) QueryExchangeListInUse(d *dto.BusGroupAccountInfoGetReq, p *actions.DataPermission, list *[]dto.CexExchangeListResp) error {
	var accountInfo models.BusExchangeAccountInfo

	err := e.Orm.Model(&accountInfo).
		Scopes(
			actions.Permission(accountInfo.TableName(), p), // 数据权限处理
		).
		Distinct("exchange_type").
		Find(list).Error

	e.Log.Infof("QueryExchangeListInUse %+v\r\n", list)

	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// GetPortfolioUnwindingInfo  实时获取账户组的借贷以及余额信息
func (e *BusExchangeAccountInfo) GetPortfolioUnwindingInfo(d *dto.ProtfolioUnwindingInfoReq, p *actions.DataPermission, resp *dto.ProtfolioUnwindingInfoResp) error {
	var exchangeAccountInfo models.BusExchangeAccountInfo
	var dexWalletInfo models.BusDexWallet

	cexAccountId := d.CexAccountId
	dexWalletId := d.DexWalletId

	err := e.Orm.Model(&exchangeAccountInfo).
		Where("id = ?", cexAccountId).
		First(&exchangeAccountInfo).Error
	if err != nil {
		e.Log.Errorf("获取cex账号失败:%s \r\n", err)
		return err
	}

	err = e.Orm.Model(&dexWalletInfo).
		Where("id = ?", dexWalletId).
		First(&dexWalletInfo).Error
	if err != nil {
		e.Log.Errorf("获取dex钱包失败:%s \r\n", err)
		return err
	}

	var masterCexAccount models.BusExchangeAccountInfo
	var secretConfig *waterLevelPb.SecretKey
	if exchangeAccountInfo.MasterAccountId == 0 {
		secretConfig, err = generateSecretConfig(dexWalletInfo, exchangeAccountInfo, models.BusExchangeAccountInfo{})
		if err != nil {
			e.Log.Errorf("生成secret config 参数失败:%s \r\n", err)
			return err
		}
	} else {
		err = e.Orm.Model(&exchangeAccountInfo).
			Where("id = ?", exchangeAccountInfo.MasterAccountId).
			First(&masterCexAccount).Error
		if err != nil {
			e.Log.Errorf("获取主账号失败:%s \r\n", err)
			return err
		}
		secretConfig, err = generateSecretConfig(dexWalletInfo, exchangeAccountInfo, masterCexAccount)
		if err != nil {
			e.Log.Errorf("生成secret config 参数失败:%s \r\n", err)
			return err
		}
	}

	exchangeType := exchangeAccountInfo.ExchangeType
	var exchangeTypeEnum waterLevelPb.ExchangeType
	if exchangeType == global.EXCHANGE_TYPE_BINANCE {
		exchangeTypeEnum = waterLevelPb.ExchangeType_Binance
	} else if exchangeType == global.EXCHANGE_TYPE_GATEIO {
		exchangeTypeEnum = waterLevelPb.ExchangeType_Gate
	} else {
		e.Log.Errorf("不支持的交易所类型:%s \r\n", exchangeType)
		return errors.New("不支持的交易所类型")
	}

	req := &waterLevelPb.PortfolioUnwindingRequest{
		TokenAddress: &d.TokenAddress,
		TokenName:    d.TokenName,
		SecretKey:    secretConfig,
		ExchangeType: exchangeTypeEnum,
	}

	grpcResp, err := client.GetPortfolioUnwindingInfo(req)
	if err != nil {
		e.Log.Errorf("获取实时借贷信息失败:%s \r\n", err)
		return err
	}
	resp.TokenName = grpcResp.TokenName
	resp.WalletBalance = grpcResp.WalletBalance
	resp.TraderAccountBorrowed = grpcResp.Borrowed
	resp.TraderAccountMarginBalance = grpcResp.TraderAccountMarginBalance
	resp.TraderAccountSpotBalance = grpcResp.TraderAccountSpotBalance
	resp.MasterAccountSpotBalance = grpcResp.MasterAccountSpotBalance

	return nil
}

// PortfolioUnwinding  资金归拢提交
func (e *BusExchangeAccountInfo) PortfolioUnwinding(d *dto.ProtfolioUnwindingInfoReq, p *actions.DataPermission) error {
	var exchangeAccountInfo models.BusExchangeAccountInfo
	var dexWalletInfo models.BusDexWallet

	cexAccountId := d.CexAccountId
	dexWalletId := d.DexWalletId

	err := e.Orm.Model(&exchangeAccountInfo).
		Where("id = ?", cexAccountId).
		First(&exchangeAccountInfo).Error
	if err != nil {
		e.Log.Errorf("获取cex账号失败:%s \r\n", err)
		return err
	}

	err = e.Orm.Model(&dexWalletInfo).
		Where("id = ?", dexWalletId).
		First(&dexWalletInfo).Error
	if err != nil {
		e.Log.Errorf("获取dex钱包失败:%s \r\n", err)
		return err
	}

	var masterCexAccount models.BusExchangeAccountInfo
	var secretConfig *waterLevelPb.SecretKey
	if exchangeAccountInfo.MasterAccountId == 0 {
		secretConfig, err = generateSecretConfig(dexWalletInfo, exchangeAccountInfo, models.BusExchangeAccountInfo{})
		if err != nil {
			e.Log.Errorf("生成secret config 参数失败:%s \r\n", err)
			return err
		}
	} else {
		err = e.Orm.Model(&exchangeAccountInfo).
			Where("id = ?", exchangeAccountInfo.MasterAccountId).
			First(&masterCexAccount).Error
		if err != nil {
			e.Log.Errorf("获取主账号失败:%s \r\n", err)
			return err
		}
		secretConfig, err = generateSecretConfig(dexWalletInfo, exchangeAccountInfo, masterCexAccount)
		if err != nil {
			e.Log.Errorf("生成secret config 参数失败:%s \r\n", err)
			return err
		}
	}

	exchangeType := exchangeAccountInfo.ExchangeType
	var exchangeTypeEnum waterLevelPb.ExchangeType
	if exchangeType == global.EXCHANGE_TYPE_BINANCE {
		exchangeTypeEnum = waterLevelPb.ExchangeType_Binance
	} else if exchangeType == global.EXCHANGE_TYPE_GATEIO {
		exchangeTypeEnum = waterLevelPb.ExchangeType_Gate
	} else {
		e.Log.Errorf("不支持的交易所类型:%s \r\n", exchangeType)
		return errors.New("不支持的交易所类型")
	}

	req := &waterLevelPb.PortfolioUnwindingRequest{
		TokenAddress: &d.TokenAddress,
		TokenName:    d.TokenName,
		SecretKey:    secretConfig,
		ExchangeType: exchangeTypeEnum,
	}

	err = client.PortfolioUnwinding(req)
	if err != nil {
		e.Log.Errorf("资金归拢失败:%s \r\n", err)
		return err
	}

	return nil
}
