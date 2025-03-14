package models

import (
	"quanta-admin/common/models"
)

type BusDexCexTriangularArbitrageDexWallet struct {
	models.Model

	WalletName          string `json:"walletName" gorm:"type:varchar(255);not null"`
	WalletAddress       string `json:"walletAddress" gorm:"type:varchar(255);not null"`
	EncryptedPrivateKey string `json:"-" gorm:"type:text;not null"`
}

func (BusDexCexTriangularArbitrageDexWallet) TableName() string {
	return "bus_dex_cex_triangular_arbitrage_dex_wallet"
}

func (e *BusDexCexTriangularArbitrageDexWallet) GetId() interface{} {
	return e.Id
}
