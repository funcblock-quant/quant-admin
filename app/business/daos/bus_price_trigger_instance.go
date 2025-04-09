package daos

import (
	"quanta-admin/app/business/models"
	"time"

	"gorm.io/gorm"
)

type BusPriceTriggerInstanceDAO struct {
	Db *gorm.DB
}

func (dao *BusPriceTriggerInstanceDAO) GetInstancesListByIds(ids []string, instances *[]models.BusPriceTriggerStrategyInstance) error {
	return dao.Db.Model(&models.BusPriceTriggerStrategyInstance{}).Where("id in ?", ids).Find(instances).Error
}

func (dao *BusPriceTriggerInstanceDAO) GetInstancesList(instances *[]models.BusPriceTriggerStrategyInstance) error {
	return dao.Db.Model(&models.BusPriceTriggerStrategyInstance{}).Find(instances).Error
}

func (dao *BusPriceTriggerInstanceDAO) ExpireInstanceWithIds(ids []string) error {
	return dao.Db.Model(&models.BusPriceTriggerStrategyInstance{}).
		Where("id in ?", ids).
		Update("status", "expired").
		Error
}

func (dao *BusPriceTriggerInstanceDAO) GetStopProfitTradingRecord(instanceId int, startTime time.Time, trades *[]models.BusPriceMonitorForOptionHedging) error {
	return dao.Db.Model(&models.BusPriceMonitorForOptionHedging{}).
		Where("strategy_instance_id = ? AND extra IS NULL AND created_at > ? AND pnl > 0", instanceId, startTime).
		Find(trades).Error
}

func (dao *BusPriceTriggerInstanceDAO) UpdateInstancePaused(instanceId int) error {
	return dao.Db.Model(&models.BusPriceMonitorForOptionHedging{}).
		Where("id = ?", instanceId).
		Update("status", "paused").
		Error
}
