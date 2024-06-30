package ports

import (
	"time"

	"go.dataflow.ru/service-sales/internal/app/domain"
)

type Storage interface {
	AddSale(sale *domain.Sale) error
	GetSales() ([]*domain.Sale, error)
	GetTotal(storeID string, startDate, endDate time.Time) (float64, error)
}
