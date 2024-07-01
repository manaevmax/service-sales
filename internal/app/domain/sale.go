package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Sale struct {
	ProductID    string
	StoreID      string
	QuantitySold int64
	SalePrice    decimal.Decimal
	SaleDate     time.Time
}
