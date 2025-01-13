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
	return dao.Db.Model(&models.BusStrategyInstance{}).Where("id = ?", request.Id).First(instance).Error
}
