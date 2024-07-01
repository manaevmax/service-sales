package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.dataflow.ru/service-sales/pkg/logger"

	"go.dataflow.ru/service-sales/internal/app/domain"
	"go.dataflow.ru/service-sales/internal/app/ports"
)

func TestService_AddSale(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	successSale := domain.Sale{
		ProductID:    "product_100",
		StoreID:      "store_1",
		QuantitySold: 10,
		SalePrice:    decimal.NewFromFloat(199),
		SaleDate:     time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
	}

	testCases := []struct {
		name    string
		sale    domain.Sale
		storage func() ports.SalesStorage
		err     error
	}{
		{
			name: "success",
			sale: successSale,
			storage: func() ports.SalesStorage {
				storage := NewMockSalesStorage(ctrl)
				storage.EXPECT().AddSale(&successSale).Return()

				return storage
			},
		},
		{
			name: "negative quantity",
			sale: domain.Sale{
				ProductID:    "product_100",
				StoreID:      "store_1",
				QuantitySold: -10,
				SalePrice:    decimal.NewFromFloat(199),
				SaleDate:     time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
			},
			storage: func() ports.SalesStorage {
				return NewMockSalesStorage(ctrl)
			},
			err: fmt.Errorf("invalid quantity"),
		},
		{
			name: "negative price",
			sale: domain.Sale{
				ProductID:    "product_100",
				StoreID:      "store_1",
				QuantitySold: 10,
				SalePrice:    decimal.NewFromFloat(-199),
				SaleDate:     time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
			},
			storage: func() ports.SalesStorage {
				return NewMockSalesStorage(ctrl)
			},
			err: fmt.Errorf("invalid sale price"),
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			saleService := NewSaleService(tt.storage(), logger.NoOpLogger())
			err := saleService.AddSale(&tt.sale)

			if tt.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tt.err, err.Error())
			}
		})
	}
}

func TestService_GetSales(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	storage := NewMockSalesStorage(ctrl)

	sales := []*domain.Sale{
		{
			ProductID:    "product_100",
			StoreID:      "store_1",
			QuantitySold: 10,
			SalePrice:    decimal.NewFromFloat(199),
			SaleDate:     time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
		},
		{
			ProductID:    "product_101",
			StoreID:      "store_1",
			QuantitySold: 10,
			SalePrice:    decimal.NewFromFloat(199),
			SaleDate:     time.Date(2024, 6, 20, 10, 10, 0, 0, time.UTC),
		},
	}

	storage.EXPECT().GetSales().Return(sales)

	saleService := NewSaleService(storage, logger.NoOpLogger())
	actualSales := saleService.GetSales()
	assert.Equal(t, sales, actualSales)
}

func TestService_GetTotalSum(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	storage := NewMockSalesStorage(ctrl)

	storeID := "store_1"
	dateFrom := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	total := decimal.NewFromFloat(1001.9)

	storage.EXPECT().GetTotalSum(storeID, dateFrom, dateTo).Return(total)

	saleService := NewSaleService(storage, logger.NoOpLogger())
	actualTotal := saleService.GetTotalSum(storeID, dateFrom, dateTo)
	assert.Equal(t, total, actualTotal)
}
