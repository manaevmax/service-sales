package services

import (
	"go.dataflow.ru/service-sales/internal/app/ports"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"go.dataflow.ru/service-sales/internal/app/domain"
)

func TestService_AddSale(t *testing.T) {
	ctrl := gomock.NewController(t)
	storage := NewMockStorage(ctrl)

	sale := domain.Sale{
		ProductID:    1000000,
		StoreID:      2,
		QuantitySold: 10,
		SalePrice:    199,
		SaleDate:     time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
	}

	testCases := []struct {
		name    string
		sale    domain.Sale
		storage func() ports.Storage
		err     error
	}{
		{
			name: "success",
			sale: domain.Sale{
				ProductID:    1000000,
				StoreID:      2,
				QuantitySold: 10,
				SalePrice:    199,
				SaleDate:     time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
			},
			storage: func() ports.Storage {
				storage.EXPECT().AddSale(&sale).Return(nil)
				return storage
			},
		},
		{
			name: "success",
			sale: domain.Sale{
				ProductID:    1000000,
				StoreID:      2,
				QuantitySold: 10,
				SalePrice:    199,
				SaleDate:     time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
			},
			storage: func() ports.Storage {
				storage.EXPECT().AddSale(&sale).Return(nil)
				return storage
			},
		},
	}

	for _, tt := range testCases {
		saleService := NewSaleService(tt.storage())
		err := saleService.AddSale(&sale)

		if tt.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, tt.err, err.Error())
		}
	}
}

func TestService_GetSales(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockStorage(ctrl)

	sales := []*domain.Sale{
		{
			ProductID:    1000000,
			StoreID:      2,
			QuantitySold: 10,
			SalePrice:    199,
			SaleDate:     time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
		},
		{
			ProductID:    1000001,
			StoreID:      2,
			QuantitySold: 10,
			SalePrice:    199,
			SaleDate:     time.Date(2024, 6, 20, 10, 10, 0, 0, time.UTC),
		},
	}

	repo.EXPECT().GetSales().Return(sales, nil)

	saleService := NewSaleService(repo)
	actualSales, err := saleService.GetSales()
	assert.NoError(t, err)
	assert.Equal(t, sales, actualSales)
}

func TestService_GetTotalForStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockStorage(ctrl)

	storeID := int64(1)
	dateFrom := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	total := 1001.9

	repo.EXPECT().GetTotal(storeID, dateFrom, dateTo).Return(total, nil)

	saleService := NewSaleService(repo)
	actualTotal, err := saleService.GetTotal(storeID, dateFrom, dateTo)
	assert.NoError(t, err)
	assert.Equal(t, total, actualTotal)
}
