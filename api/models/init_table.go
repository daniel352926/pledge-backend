package models

import "pledge-backend/db"

func InitTable() {
	db.Mysql.AutoMigrate(&MultiSign{},
		&TokenInfo{},
		&TokenList{},
		&PoolData{},
		&PoolBases{},
		&Tx{},
		&Receipt{},
		&Block{},
	)
}
