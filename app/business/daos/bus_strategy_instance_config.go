package daos

import (
	"gorm.io/gorm"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
)

type BusStrategyInstanceConfigDAO struct {
	Db *gorm.DB
}

func (dao *BusStrategyInstanceConfigDAO) GetInstanceConfigs(request *dto.BusStrategyInstanceConfigGetByInstanceIdReq, list *[]models.BusStrategyInstanceConfig) error {
	return dao.Db.Model(&models.BusStrategyInstanceConfig{}).Where("strategy_instance_id = ?", request.StrategyInstanceId).Find(list).Error
}
