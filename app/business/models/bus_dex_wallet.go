package models

import (
	"quanta-admin/common/models"
)

type BusDexWallet struct {
	models.Model

	WalletName          string `json:"walletName" gorm:"type:varchar(255);not null"`
	WalletAddress       string `json:"walletAddress" gorm:"type:varchar(255);not null"`
	EncryptedPrivateKey string `json:"-" gorm:"type:text;not null"`
	Blockchain          string `json:"blockchain" gorm:"type:varchar(32);not null"`
}

func (BusDexWallet) TableName() string {
	return "bus_dex_wallet"
}

func (e *BusDexWallet) GetId() interface{} {
	return e.Id
}
