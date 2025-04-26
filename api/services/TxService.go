package services

import (
	"pledge-backend/api/models"
)

type TxService struct {
}

func NewTx() *TxService {
	return &TxService{}
}

func (s *TxService) Save(tx *models.Tx) error {
	return models.NewTx().Save(tx)
}

func (s *TxService) GetByHash(hash string) (*models.Tx, error) {
	return models.NewTx().GetByHash(hash)
}

func (s *TxService) GetByHashes(hashes []string) ([]models.Tx, error) {
	return models.NewTx().GetByHashes(hashes)
}
