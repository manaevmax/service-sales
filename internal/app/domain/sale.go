package domain

import "time"

type Sale struct {
	ProductID    string
	StoreID      string
	QuantitySold int64
	SalePrice    float64
	SaleDate     time.Time
}
