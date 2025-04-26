package models

import (
	"errors"
	"gorm.io/gorm"
	"pledge-backend/db"
	"pledge-backend/log"
	"time"
)

type Block struct {
	ID               uint64 `gorm:"primaryKey;autoIncrement;not null"`
	Number           string `gorm:"type:varchar(255);not null"`
	Hash             string `gorm:"type:varchar(255);not null"`
	ParentHash       string `gorm:"type:varchar(255);not null"`
	Nonce            string `gorm:"type:varchar(255);not null"`
	Sha3Uncles       string `gorm:"type:varchar(255);not null"`
	LogsBloom        string `gorm:"type:text;not null"`
	TransactionsRoot string `gorm:"type:varchar(255);not null"`
	StateRoot        string `gorm:"type:varchar(255);not null"`
	ReceiptsRoot     string `gorm:"type:varchar(255);not null"`
	Miner            string `gorm:"type:varchar(255);not null"`
	Difficulty       string `gorm:"type:varchar(255);not null"`
	TotalDifficulty  string `gorm:"type:varchar(255);not null"`
	ExtraData        string `gorm:"type:text;not null"`
	Size             string `gorm:"type:varchar(255);not null"`
	GasLimit         string `gorm:"type:varchar(255);not null"`
	GasUsed          string `gorm:"type:varchar(255);not null"`
	Timestamp        string `gorm:"type:varchar(255);not null"`
	Transactions     string `gorm:"type:json;not null"` // 需要 import "gorm.io/datatypes"
	Uncles           string `gorm:"type:json;not null"`
	CreatedAt        time.Time
}

func NewBlock() *Block {
	return &Block{}
}

func (r *Block) Save(block *Block) error {
	exists := Tx{}

	err := db.Mysql.Table("block").Where("hash=?", block.Hash).First(&exists).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Mysql.Table("block").Create(block).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err
			}
		} else {
			return errors.New("block record select err " + err.Error())
		}
	}

	return nil
}

func (r *Block) GetByNum(txHash string) (*Block, error) {
	receipt := Block{}

	err := db.Mysql.Table("block").Where("number=?", txHash).First(&receipt).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error(err.Error())
		return nil, err
	}

	return &receipt, nil
}
