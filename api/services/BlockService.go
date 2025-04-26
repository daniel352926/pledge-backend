package services

import "pledge-backend/api/models"

type BlockService struct {
}

func NewBlockService() *BlockService {
	return &BlockService{}
}

func getByHash() {

}
func (b *BlockService) Save(block *models.Block) error {
	return models.NewBlock().Save(block)
}

func (b *BlockService) GetByNum(txHash string) (*models.Block, error) {
	return models.NewBlock().GetByNum(txHash)
}
