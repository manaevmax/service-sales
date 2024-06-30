package storage

import (
	"go.dataflow.ru/service-sales/internal/app/domain"
	"sync"
	"time"
)

type SalesStorage struct {
	sales []*domain.Sale
	mu    sync.Mutex
}

func New() *SalesStorage {
	return &SalesStorage{
		sales: []*domain.Sale{},
	}
}

func (repo *SalesStorage) AddSale(sale *domain.Sale) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.sales = append(repo.sales, sale)
	return nil
}

func (repo *SalesStorage) GetSales() ([]*domain.Sale, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.sales, nil
}

func (repo *SalesStorage) GetTotal(storeID string, startDate, endDate time.Time) (float64, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var total float64
	for _, sale := range repo.sales {
		if sale.StoreID == storeID && sale.SaleDate.After(startDate) && sale.SaleDate.Before(endDate) {
			total += float64(sale.QuantitySold) * sale.SalePrice
		}
	}
	return total, nil
}
