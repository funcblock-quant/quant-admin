package dto

import (

	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusDexWalletGetPageReq struct {
	dto.Pagination     `search:"-"`
    WalletName string `form:"walletName"  search:"type:exact;column:wallet_name;table:bus_dex_wallet" comment:"钱包名称"`
    WalletAddress string `form:"walletAddress"  search:"type:exact;column:wallet_address;table:bus_dex_wallet" comment:"钱包地址"`
    Blockchain string `form:"blockchain"  search:"type:exact;column:blockchain;table:bus_dex_wallet" comment:"链网络"`
    BusDexWalletOrder
}

type BusDexWalletOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:bus_dex_wallet"`
    WalletName string `form:"walletNameOrder"  search:"type:order;column:wallet_name;table:bus_dex_wallet"`
    WalletAddress string `form:"walletAddressOrder"  search:"type:order;column:wallet_address;table:bus_dex_wallet"`
    EncryptedPrivateKey string `form:"encryptedPrivateKeyOrder"  search:"type:order;column:encrypted_private_key;table:bus_dex_wallet"`
    Blockchain string `form:"blockchainOrder"  search:"type:order;column:blockchain;table:bus_dex_wallet"`
    
}

func (m *BusDexWalletGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusDexWalletInsertReq struct {
    Id int `json:"-" comment:"主键 ID"` // 主键 ID
    WalletName string `json:"walletName" comment:"钱包名称"`
    WalletAddress string `json:"walletAddress" comment:"钱包地址"`
    EncryptedPrivateKey string `json:"encryptedPrivateKey" comment:"加密后的密钥"`
    Blockchain string `json:"blockchain" comment:"链网络"`
    common.ControlBy
}

func (s *BusDexWalletInsertReq) Generate(model *models.BusDexWallet)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.WalletName = s.WalletName
    model.WalletAddress = s.WalletAddress
    model.EncryptedPrivateKey = s.EncryptedPrivateKey
    model.Blockchain = s.Blockchain
}

func (s *BusDexWalletInsertReq) GetId() interface{} {
	return s.Id
}

type BusDexWalletUpdateReq struct {
    Id int `uri:"id" comment:"主键 ID"` // 主键 ID
    WalletName string `json:"walletName" comment:"钱包名称"`
    WalletAddress string `json:"walletAddress" comment:"钱包地址"`
    EncryptedPrivateKey string `json:"encryptedPrivateKey" comment:"加密后的密钥"`
    Blockchain string `json:"blockchain" comment:"链网络"`
    common.ControlBy
}

func (s *BusDexWalletUpdateReq) Generate(model *models.BusDexWallet)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.WalletName = s.WalletName
    model.WalletAddress = s.WalletAddress
    model.EncryptedPrivateKey = s.EncryptedPrivateKey
    model.Blockchain = s.Blockchain
}

func (s *BusDexWalletUpdateReq) GetId() interface{} {
	return s.Id
}

// BusDexWalletGetReq 功能获取请求参数
type BusDexWalletGetReq struct {
     Id int `uri:"id"`
}
func (s *BusDexWalletGetReq) GetId() interface{} {
	return s.Id
}

// BusDexWalletDeleteReq 功能删除请求参数
type BusDexWalletDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusDexWalletDeleteReq) GetId() interface{} {
	return s.Ids
}
