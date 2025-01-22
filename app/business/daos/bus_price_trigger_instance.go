package daos

import (
	"gorm.io/gorm"
	"quanta-admin/app/business/models"
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
