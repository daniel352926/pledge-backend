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
	Number           uint64 `gorm:"type:bigint;not null"`
	Hash             string `gorm:"type:varchar(255);not null"`
	ParentHash       string `gorm:"type:varchar(255);not null"`
	Nonce            uint64 `gorm:"type:bigint;not null"`
	Sha3Uncles       string `gorm:"type:varchar(255);not null"`
	LogsBloom        string `gorm:"type:text;not null"`
	TransactionsRoot string `gorm:"type:varchar(255);not null"`
	StateRoot        string `gorm:"type:varchar(255);not null"`
	ReceiptsRoot     string `gorm:"type:varchar(255);not null"`
	Miner            string `gorm:"type:varchar(255);not null"`
	Difficulty       string `gorm:"type:varchar(255);not null"`
	ExtraData        string `gorm:"type:text;not null"`
	Size             string `gorm:"type:varchar(255);not null"`
	GasLimit         string `gorm:"type:varchar(255);not null"`
	GasUsed          string `gorm:"type:varchar(255);not null"`
	Timestamp        uint64 `gorm:"type:int;not null"`
	Transactions     string `gorm:"type:json;not null"`
	Uncles           string `gorm:"type:json;not null"`
	CreatedAt        time.Time
}

func NewBlock() *Block {
	return &Block{}
}

func (r *Block) Save(block *Block) error {
	exists, err := r.GetByNum(block.Number)
	if err != nil {
		return errors.New("block record select err " + err.Error())
	}

	if exists == nil {
		err = db.Mysql.Table("block").Create(block).Debug().Error
		if err != nil {
			log.Logger.Error(err.Error())
			return err
		}
	}

	return nil
}

func (r *Block) GetByNum(blockNumber uint64) (*Block, error) {
	block := Block{}

	err := db.Mysql.Table("block").Where("number=?", blockNumber).First(&block).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Logger.Error(err.Error())
		return nil, err
	}

	return &block, nil
}
