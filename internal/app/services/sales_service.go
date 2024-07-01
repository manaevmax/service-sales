package services

//go:generate mockgen -package $GOPACKAGE -source ../ports/sales_storage.go -destination mocks.go

import (
	"fmt"
	"time"

	"go.dataflow.ru/service-sales/internal/app/domain"
	"go.dataflow.ru/service-sales/internal/app/ports"
	"go.dataflow.ru/service-sales/pkg/logger"
)

type SalesService struct {
	storage ports.SalesStorage
	logger  *logger.Logger
}

func NewSaleService(storage ports.SalesStorage, logger *logger.Logger) *SalesService {
	return &SalesService{
		storage: storage,
		logger:  logger,
	}
}

func (s *SalesService) AddSale(sale *domain.Sale) error {
	if sale.SalePrice < 0 {
		return fmt.Errorf("invalid sale price")
	}

	if sale.QuantitySold < 0 {
		return fmt.Errorf("invalid quantity")
	}

	s.storage.AddSale(sale)

	return nil
}

func (s *SalesService) GetSales() []*domain.Sale {
	return s.storage.GetSales()
}

func (s *SalesService) GetTotalSum(storeID string, startDate, endDate time.Time) float64 {
	return s.storage.GetTotalSum(storeID, startDate, endDate)
}
