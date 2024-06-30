package services

//go:generate mockgen -package $GOPACKAGE -source ../ports/storage.go -destination mocks.go

import (
	"fmt"
	"time"

	"go.dataflow.ru/service-sales/internal/app/domain"
	"go.dataflow.ru/service-sales/internal/app/ports"
)

type SaleService struct {
	storage ports.Storage
}

func NewSaleService(repo ports.Storage) *SaleService {
	return &SaleService{
		storage: repo,
	}
}

func (s *SaleService) AddSale(sale *domain.Sale) error {
	if sale.SalePrice < 0 {
		return fmt.Errorf("invalid sale price")
	}

	if sale.QuantitySold < 0 {
		return fmt.Errorf("invalid quantity")
	}

	return s.storage.AddSale(sale)
}

func (s *SaleService) GetSales() ([]*domain.Sale, error) {
	return s.storage.GetSales()
}

func (s *SaleService) GetTotal(storeID string, startDate, endDate time.Time) (float64, error) {
	return s.storage.GetTotal(storeID, startDate, endDate)
}
