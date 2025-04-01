package models

import (
	"quanta-admin/common/models"
)

type BusDexWallet struct {
	models.Model
	WalletName          string `json:"walletName" gorm:"type:varchar(255);comment:钱包名称"`
	WalletAddress       string `json:"walletAddress" gorm:"type:varchar(255);comment:钱包地址"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey" gorm:"type:varchar(512);comment:加密后的密钥"`
	Blockchain          string `json:"blockchain" gorm:"type:varchar(32);comment:链网络"`
}

func (BusDexWallet) TableName() string {
	return "bus_dex_wallet"
}

func (e *BusDexWallet) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusDexWallet) GetId() interface{} {
	return e.Id
}

// 实现 SetCreateBy 方法
func (e *BusDexWallet) SetCreateBy(userId int) {
}

// 实现 SetCreateBy
func (e *BusDexWallet) SetUpdateBy(userId int) {
}
