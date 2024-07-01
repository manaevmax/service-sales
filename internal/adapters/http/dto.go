package http

import (
	"fmt"
	"time"
)

type SaleDto struct {
	ProductID    string  `json:"product_id"`
	StoreID      string  `json:"store_id"`
	QuantitySold int64   `json:"quantity_sold"`
	SalePrice    float64 `json:"sale_price"`
	SaleDate     string  `json:"sale_date"`
}

func (r *SaleDto) Validate() error {
	if r.ProductID == "" {
		return fmt.Errorf("productID not defined")
	}

	if r.StoreID == "" {
		return fmt.Errorf("storeID not defined")
	}

	if r.QuantitySold <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	if r.SalePrice <= 0 {
		return fmt.Errorf("price must be positive")
	}

	if _, err := time.Parse(time.RFC3339, r.SaleDate); err != nil {
		return fmt.Errorf("date must be in RFC3339 format")
	}

	return nil
}

type AddSaleResponse struct {
	Status string `json:"status"`
}

type CalculateTotalSumRequest struct {
	Operation string `json:"operation"`
	StoreID   string `json:"store_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (r *CalculateTotalSumRequest) Validate() error {
	if r.Operation != "total_sales" {
		return fmt.Errorf("unknown operation")
	}

	return nil
}

type CalculateTotalSumResponse struct {
	StoreID    string  `json:"store_id"`
	TotalSales float64 `json:"total_sales"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
}
