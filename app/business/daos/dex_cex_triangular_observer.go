package daos

import (
	"gorm.io/gorm"
	"quanta-admin/app/business/models"
)

type BusDexCexTriangularObserverDAO struct {
	Db *gorm.DB
}

func (dao *BusDexCexTriangularObserverDAO) GetObserverList(observers *[]models.BusDexCexTriangularObserver) error {
	return dao.Db.Model(&models.BusDexCexTriangularObserver{}).Find(observers).Error
}

func (dao *BusDexCexTriangularObserverDAO) UpdateObserverWithNewId(observerId string, id int) error {
	return dao.Db.Model(&models.BusDexCexTriangularObserver{}).
		Where("id = ?", id).
		Update("instance_id", observerId).
		Error
}

func (dao *BusDexCexTriangularObserverDAO) UpdateObserverWithTradingStatus(id int, isTrading bool) error {
	return dao.Db.Model(&models.BusDexCexTriangularObserver{}).
		Where("id = ?", id).
		Update("is_trading", isTrading).
		Error
}
