package services

import "pledge-backend/api/models"

type ReceiptService struct {
}

func NewReceipt() *ReceiptService {
	return &ReceiptService{}
}

func (r *ReceiptService) Save(receipt *models.Receipt) error {
	return models.NewReceipt().Save(receipt)
}

func (r *ReceiptService) GetByHash(txHash string) (*models.Receipt, error) {
	return models.NewReceipt().GetByHash(txHash)
}
