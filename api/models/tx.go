package models

import (
	"errors"
	"gorm.io/gorm"
	"pledge-backend/db"
	"pledge-backend/log"
)

type Tx struct {
	Id                   int    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	Hash                 string `json:"hash" gorm:"column:tx_hash;type:char(66);not null"`
	Type                 uint8  `json:"type" gorm:"column:tx_type;type:tinyint"`
	Nonce                uint64 `json:"nonce" gorm:"column:nonce;type:bigint"`
	GasPrice             string `json:"gasPrice" gorm:"column:gas_price;type:varchar(66);null"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas" gorm:"column:max_priority_fee_per_gas;type:varchar(66)"`
	MaxFeePerGas         string `json:"maxFeePerGas" gorm:"column:max_fee_per_gas;type:varchar(66)"`
	Gas                  string `json:"gas" gorm:"column:gas_limit;type:varchar(20)"`
	Value                string `json:"value" gorm:"column:value;type:varchar(66)"`
	Input                string `json:"input" gorm:"column:input_data;type:longtext"`
	V                    string `json:"v" gorm:"column:v;type:varchar(10)"`
	R                    string `json:"r" gorm:"column:r;type:char(66)"`
	S                    string `json:"s" gorm:"column:s;type:char(66)"`
	To                   string `json:"to" gorm:"column:to_address;type:char(42)"`
	ChainId              string `json:"chainId" gorm:"column:chain_id;type:varchar(20)"`
	AccessList           string `json:"accessList" gorm:"column:access_list;type:json"`
}

func NewTx() *Tx {
	return &Tx{}
}

func (t Tx) Save(tx *Tx) error {
	exists := Tx{}

	err := db.Mysql.Table("transaction").Where("tx_hash=?", tx.Hash).First(&exists).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Mysql.Table("transaction").Create(tx).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err
			}
		} else {
			return errors.New("transaction record select err " + err.Error())
		}
	}

	return nil
}

func (t Tx) GetByHash(hash string) (*Tx, error) {
	res := Tx{}
	err := db.Mysql.Table("transaction").Where("tx_hash=?", hash).First(&res).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error(err.Error())
		return nil, err
	}
	return &res, nil
}

func (t Tx) GetByHashes(hashes []string) ([]Tx, error) {
	var res []Tx
	err := db.Mysql.Table("transaction").Where("tx_hash in ?", hashes).Find(&res).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error(err.Error())
		return nil, err
	}
	return res, nil
}
