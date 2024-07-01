package ports

import (
	"time"

	"go.dataflow.ru/service-sales/internal/app/domain"
)

type SalesStorage interface {
	AddSale(sale *domain.Sale)
	GetSales() []*domain.Sale
	GetTotalSum(storeID string, startDate, endDate time.Time) float64
}
