package daos

import (
	"gorm.io/gorm"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
)

type BusStrategyInstanceDAO struct {
	Db *gorm.DB
}

func (dao *BusStrategyInstanceDAO) GetInstanceStartStopFlag(request *dto.BusStrategyInstanceGetReq, instance *models.BusStrategyInstance) error {
	return dao.Db.Model(&models.BusStrategyInstance{}).
		Where("id = ?", request.Id).
		First(instance).Error
}

// GetRegisteredInstanceList 获取所有注册的策略列表
func (dao *BusStrategyInstanceDAO) GetRegisteredInstanceList(results *[]models.BusStrategyBaseInfo) error {
	return dao.Db.Model(&models.BusStrategyBaseInfo{}).
		Where("status = 1 and deleted_at is null").
		Find(results).Error
}

func (dao *BusStrategyInstanceDAO) GetRunningInstanceByStrategyId(strategyId int, instances *[]models.BusStrategyInstance) error {
	return dao.Db.Model(&models.BusStrategyInstance{}).
		Where("strategy_id = ? and status = ? and deleted_at is null", strategyId, 1). // status = 1 表示已开启
		Find(instances).Error
}

func (dao *BusStrategyInstanceDAO) GetInstanceConfigByInstanceId(instanceId int, instanceConfig *models.BusStrategyInstanceConfig) error {
	return dao.Db.Model(&models.BusStrategyInstanceConfig{}).
		Where("strategy_instance_id = ? and deleted_at is null", instanceId).
		First(instanceConfig).Error

}
