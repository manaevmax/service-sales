package http

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"

	"go.dataflow.ru/service-sales/internal/app/domain"
	"go.dataflow.ru/service-sales/internal/app/services"
)

type SalesHandler struct {
	SaleService *services.SaleService
}

func NewSalesHandler(service *services.SaleService) *SalesHandler {
	return &SalesHandler{SaleService: service}
}

func (h *SalesHandler) AddSale(c *fiber.Ctx) error {
	var req SaleDto

	err := c.BodyParser(&req)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	if err = req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = h.SaleService.AddSale(convertFromDto(req))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *SalesHandler) GetSales(c *fiber.Ctx) error {
	sales, err := h.SaleService.GetSales()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(sales)
}

func (h *SalesHandler) CalculateTotalSales(c *fiber.Ctx) error {
	var req TotalSalesRequest

	err := c.BodyParser(&req)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	if err = req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	totalSales, err := h.SaleService.GetTotal(req.StoreID, startDate, endDate)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"store_id":    request.StoreID,
		"total_sales": totalSales,
		"start_date":  request.StartDate,
		"end_date":    request.EndDate,
	}
	json.NewEncoder(w).Encode(response)

	return
}

func convertFromDto(s SaleDto) *domain.Sale {
	dt, _ := time.Parse(time.RFC3339, s.SaleDate)
	return &domain.Sale{
		ProductID:    s.ProductID,
		StoreID:      s.StoreID,
		QuantitySold: s.QuantitySold,
		SalePrice:    s.SalePrice,
		SaleDate:     dt,
	}
}
