package storage

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"go.dataflow.ru/service-sales/internal/app/domain"
	"go.dataflow.ru/service-sales/pkg/logger"
)

const (
	defaultIndexGranularity = 10
	missedIndex             = -1
)

// SaleExt определяет тип для продаж, расширенный кумулятивной суммой продаж.
type SaleExt struct {
	*domain.Sale
	CumulativeSum decimal.Decimal
}

// SalesStorage хранилище для работы с продажами.
type SalesStorage struct {
	salesByStore map[string][]SaleExt

	sparseIndex      map[string][]time.Time // разреженный индекс для хранения временных меток продаж
	indexGranularity int64

	mu sync.RWMutex

	logger *logger.Logger
}

// New возвращает новый экземпляр хранилища.
func New(logger *logger.Logger, opts ...Option) *SalesStorage {
	s := &SalesStorage{
		salesByStore:     make(map[string][]SaleExt),
		sparseIndex:      make(map[string][]time.Time),
		mu:               sync.RWMutex{},
		logger:           logger,
		indexGranularity: defaultIndexGranularity,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// AddSale сохраняет информацию о продаже.
func (s *SalesStorage) AddSale(sale *domain.Sale) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sales := s.salesByStore[sale.StoreID]

	// сохраняем разреженный индекс, если необходимо
	if len(sales)%defaultIndexGranularity == 0 {
		s.sparseIndex[sale.StoreID] = append(s.sparseIndex[sale.StoreID], sale.SaleDate)
	}

	// считаем кумулятивную сумму продаж для текущей продажи
	cumulativeSum := decimal.NewFromFloat(0)
	if len(sales) > 0 {
		cumulativeSum = sales[len(sales)-1].CumulativeSum
	}

	cumulativeSum = cumulativeSum.Add(decimal.NewFromInt(sale.QuantitySold).Mul(sale.SalePrice))

	// сохраняем информацию о продаже и кумулятивной сумме продаж
	s.salesByStore[sale.StoreID] = append(s.salesByStore[sale.StoreID], SaleExt{
		Sale:          sale,
		CumulativeSum: cumulativeSum,
	})
}

// GetSales возвращает данные о всех продажах.
func (s *SalesStorage) GetSales() []*domain.Sale {
	s.mu.RLock()
	defer s.mu.RUnlock()

	salesCount := 0
	for _, storeSales := range s.salesByStore {
		salesCount += len(storeSales)
	}

	sales := make([]*domain.Sale, 0, salesCount)

	for _, storeSales := range s.salesByStore {
		for _, s := range storeSales {
			sales = append(sales, s.Sale)
		}
	}

	return sales
}

// GetTotalSum возвращает сумму продаж магазина за период.
func (s *SalesStorage) GetTotalSum(storeID string, startDate, endDate time.Time) decimal.Decimal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sales := s.salesByStore[storeID]

	if startDate.After(endDate) || len(sales) == 0 {
		return decimal.NewFromFloat(0)
	}

	// обрабатываем корнер-кейсы, когда запрошенный диапазон не пересекается с периодом, за который есть продажи
	if sales[0].SaleDate.After(endDate) {
		return decimal.NewFromFloat(0)
	}

	if sales[len(sales)-1].SaleDate.Before(startDate) {
		return decimal.NewFromFloat(0)
	}

	// считаем сумму продаж за период
	firstSaleIdx := s.findSaleIdx(storeID, startDate)
	lastSaleIdx := s.findSaleIdx(storeID, endDate)

	firstSaleCumulativeSum := decimal.NewFromFloat(0)
	if firstSaleIdx > 0 {
		firstSaleCumulativeSum = sales[firstSaleIdx-1].CumulativeSum
	}

	lastSaleCumulativeSum := sales[len(sales)-1].CumulativeSum
	if lastSaleIdx < int64(len(sales)-1) && lastSaleIdx >= 0 {
		lastSaleCumulativeSum = sales[lastSaleIdx].CumulativeSum
	}

	return lastSaleCumulativeSum.Sub(firstSaleCumulativeSum)
}

// findSaleIdx возвращает индекс продажи, ближайшей справа на временной оси к заданному времени dt.
// Если продажа не найдена - возвращает missedIndex.
func (s *SalesStorage) findSaleIdx(storeID string, dt time.Time) int64 {
	sales := s.salesByStore[storeID]

	granuleHeadIdx := s.findGranuleHeadIdx(storeID, dt)

	granuleTailIdx := granuleHeadIdx + s.indexGranularity // "хвост" не входит в гранулу
	if int64(len(sales)) < granuleTailIdx {
		granuleTailIdx = int64(len(sales))
	}

	// ищем продажу перебором по всей грануле.
	for i := granuleHeadIdx; i < granuleTailIdx; i++ {
		if sales[i].SaleDate.After(dt) || sales[i].SaleDate == dt {
			return i
		}
	}

	return missedIndex
}

// findGranuleHeadIdx возвращает стартовый индекс гранулы, внутри которой находится искомая продажа.
func (s *SalesStorage) findGranuleHeadIdx(storeID string, dt time.Time) int64 {
	for i, t := range s.sparseIndex[storeID] {
		if t.After(dt) {
			if i == 0 {
				return 0
			}

			return int64(i-1) * defaultIndexGranularity
		}
	}

	// продажа с искомой временной меткой в последней грануле или не существует
	return int64(len(s.sparseIndex[storeID])-1) * defaultIndexGranularity
}

// GetTotalSumSimple возвращает сумму продаж магазина за период простым перебором (для сравнения).
func (s *SalesStorage) GetTotalSumSimple(storeID string, startDate, endDate time.Time) decimal.Decimal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sum := decimal.NewFromFloat(0)

	for _, sale := range s.salesByStore[storeID] {
		if (sale.SaleDate.After(startDate) || sale.SaleDate == startDate) && (sale.SaleDate.Before(endDate) || sale.SaleDate == endDate) {
			sum = sum.Add(decimal.NewFromInt(sale.QuantitySold).Mul(sale.SalePrice))
		}
	}

	return sum
}
