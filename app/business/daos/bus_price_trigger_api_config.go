package daos

import (
	"gorm.io/gorm"
	"quanta-admin/app/business/models"
)

type BusPriceTriggerApiConfigDAO struct {
	Db *gorm.DB
}

func (dao *BusPriceTriggerApiConfigDAO) GetApiConfigById(id int, apiConfig *models.BusPriceTriggerStrategyApikeyConfig) error {
	return dao.Db.Model(&models.BusPriceTriggerStrategyApikeyConfig{}).Where("id =", id).First(apiConfig).Error
}
