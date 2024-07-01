package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.dataflow.ru/service-sales/internal/app/domain"
	"go.dataflow.ru/service-sales/pkg/logger"
)

func TestSalesStorage_GetTotal(t *testing.T) {
	s := New(logger.NoOpLogger())

	// готовим данные для тестов
	// сумма всех продаж sum(1..98) = 4851, сумма каждой следующей продажи на 1 больше предыдущей
	// первая продажа 2024-05-01, последняя продажа 2024-08-06, шаг 1 день
	dt := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 98; i++ {
		s.AddSale(&domain.Sale{
			StoreID:      "store_1",
			ProductID:    fmt.Sprintf("product_%d", i),
			QuantitySold: 1,
			SalePrice:    float64(1 + i),
			SaleDate:     dt.AddDate(0, 0, i),
		})
	}

	testCases := []struct {
		name     string
		dateFrom time.Time
		dateTo   time.Time
		expSum   float64
	}{
		{
			name:     "dateFrom > dateTo",
			dateFrom: time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 6, 19, 0, 0, 0, 0, time.UTC),
			expSum:   0,
		},
		{
			name:     "dateFrom и dateTo левее всех продаж",
			dateFrom: time.Date(2024, 4, 29, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
			expSum:   0,
		},
		{
			name:     "dateFrom левее всех продаж, dateTo попадает на первую продажу",
			dateFrom: time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			expSum:   1,
		},
		{
			name:     "dateFrom левее всех продаж, dateTo попадает на первую продажу",
			dateFrom: time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			expSum:   1,
		},
		{
			name:     "dateFrom левее всех продаж, dateTo попадает на вторую продажу",
			dateFrom: time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 5, 2, 0, 0, 0, 0, time.UTC),
			expSum:   3,
		},
		{
			name:     "dateFrom левее всех продаж, dateTo попадает на последнюю продажу",
			dateFrom: time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
			expSum:   4851,
		},
		{
			name:     "dateFrom левее всех продаж, dateTo правее последней продажи",
			dateFrom: time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 8, 7, 0, 0, 0, 0, time.UTC),
			expSum:   4851,
		},
		{
			name:     "dateFrom на первой продаже, dateTo на первой продаже",
			dateFrom: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			expSum:   1,
		},
		{
			name:     "dateFrom на первой продаже, dateTo на третьей продаже",
			dateFrom: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC),
			expSum:   6,
		},
		{
			name:     "dateFrom на предпоследней продаже, dateTo на последней продаже",
			dateFrom: time.Date(2024, 8, 5, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
			expSum:   195,
		},
		{
			name:     "dateFrom на последней продаже, dateTo правее всех продаж",
			dateFrom: time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 8, 7, 0, 0, 0, 0, time.UTC),
			expSum:   98,
		},
		{
			name:     "dateFrom и dateTo правее всех продаж",
			dateFrom: time.Date(2024, 8, 7, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 8, 8, 0, 0, 0, 0, time.UTC),
			expSum:   0,
		},
		{
			name:     "dateFrom на первой продаже, dateTo на последней продаже",
			dateFrom: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			dateTo:   time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
			expSum:   4851,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expSum, s.GetTotalSum("store_1", tt.dateFrom, tt.dateTo))
		})
	}
}

func TestSalesStorage_GetTotal_Performance(t *testing.T) {
	t.Parallel()

	s := New(logger.NoOpLogger(), WithIndexGranularity(1000))

	dateFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom
	for i := 0; i < 10000000; i++ {
		s.AddSale(&domain.Sale{
			StoreID:      "store_1",
			ProductID:    fmt.Sprintf("product_%d", i),
			QuantitySold: 1,
			SalePrice:    float64(1 + i),
			SaleDate:     dateFrom.Add(time.Duration(i) * time.Second),
		})

		dateTo = dateTo.Add(time.Duration(i) * time.Second)
	}

	durationSimpleStart := time.Now()
	totalSimple := s.GetTotalSumSimple("store_1", dateFrom, dateTo)
	durationSimple := time.Since(durationSimpleStart)

	durationOptimizedStart := time.Now()
	totalOptimized := s.GetTotalSum("store_1", dateFrom, dateTo)
	durationOptimized := time.Since(durationOptimizedStart)

	t.Logf("performance: simple=%s, optimized=%s", durationSimple, durationOptimized)
	assert.Equal(t, totalSimple, totalOptimized)
}

func TestSalesStorage_AddGetSale(t *testing.T) {
	t.Parallel()

	s := New(logger.NoOpLogger())

	s1 := &domain.Sale{
		StoreID:      "store_1",
		ProductID:    "product_1",
		QuantitySold: 100,
		SalePrice:    9.99,
		SaleDate:     time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
	}

	s2 := &domain.Sale{
		StoreID:      "store_1",
		ProductID:    "product_2",
		QuantitySold: 100,
		SalePrice:    29.99,
		SaleDate:     time.Date(2024, 6, 2, 10, 0, 0, 0, time.UTC),
	}

	s.AddSale(s1)
	s.AddSale(s2)

	assert.Equal(t, []*domain.Sale{s1, s2}, s.GetSales())
}
