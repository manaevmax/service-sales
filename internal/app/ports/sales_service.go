package ports

import (
	"time"

	"github.com/shopspring/decimal"

	"go.dataflow.ru/service-sales/internal/app/domain"
)

type SalesService interface {
	AddSale(sale *domain.Sale) error
	GetSales() []*domain.Sale
	GetTotalSum(storeID string, startDate, endDate time.Time) decimal.Decimal
}
