package models

import (
	"errors"
	"gorm.io/gorm"
	"pledge-backend/db"
	"pledge-backend/log"
	"time"
)

type Receipt struct {
	Id                int64     `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	Hash              string    `json:"transactionHash" gorm:"column:tx_hash;type:char(66);not null"`
	Type              string    `json:"type" gorm:"column:tx_type;type:varchar(10)"`
	Status            string    `json:"status" gorm:"column:status;type:varchar(10)"`
	Root              string    `json:"root" gorm:"column:root;type:varchar(66)"`
	CumulativeGasUsed string    `json:"cumulativeGasUsed" gorm:"column:cumulative_gas_used;type:varchar(66)"`
	LogsBloom         string    `json:"logsBloom" gorm:"column:logs_bloom;type:text"`
	ContractAddress   string    `json:"contractAddress" gorm:"column:contract_address;type:char(42)"`
	GasUsed           string    `json:"gasUsed" gorm:"column:gas_used;type:varchar(66)"`
	BlockHash         string    `json:"blockHash" gorm:"column:block_hash;type:char(66)"`
	BlockNumber       string    `json:"blockNumber" gorm:"column:block_number;type:varchar(20)"`
	TransactionIndex  string    `json:"transactionIndex" gorm:"column:transaction_index;type:varchar(10)"`
	Logs              string    `json:"logs" gorm:"column:logs;type:json"` // 需要导入 "gorm.io/datatypes"
	CreatedAt         time.Time `json:"-" gorm:"column:created_at;autoCreateTime"`
}

func NewReceipt() *Receipt {
	return &Receipt{}
}

func (r *Receipt) Save(receipt *Receipt) error {
	exists := Tx{}

	err := db.Mysql.Table("receipt").Where("tx_hash=?", receipt.Hash).First(&exists).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Mysql.Table("receipt").Create(receipt).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err
			}
		} else {
			return errors.New("receipt record select err " + err.Error())
		}
	}

	return nil
}

func (r *Receipt) GetByHash(txHash string) (*Receipt, error) {
	receipt := Receipt{}

	err := db.Mysql.Table("receipt").Where("tx_hash=?", txHash).First(&receipt).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error(err.Error())
		return nil, err
	}

	return &receipt, nil
}
