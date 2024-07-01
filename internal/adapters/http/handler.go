package http

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"go.dataflow.ru/service-sales/internal/app/domain"
	"go.dataflow.ru/service-sales/internal/app/ports"
)

// SalesHandler обработчик продаж.
type SalesHandler struct {
	salesService ports.SalesService
}

// New возвращает новый экземпляр обработчика.
func New(service ports.SalesService) *SalesHandler {
	return &SalesHandler{salesService: service}
}

// AddSale обрабатывает запрос на добавление новой продажи.
func (h *SalesHandler) AddSale(c *fiber.Ctx) error {
	var req SaleDto

	err := c.BodyParser(&req)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	if err = req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = h.salesService.AddSale(convertFromDto(req))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(AddSaleResponse{Status: "success"})
}

// GetSales обрабатывает запрос получения списка всех продаж.
func (h *SalesHandler) GetSales(c *fiber.Ctx) error {
	sales := h.salesService.GetSales()

	return c.JSON(sales)
}

// CalculateTotalSum обрабатывает запрос на расчет суммы продаж для заданного магазина за период.
func (h *SalesHandler) CalculateTotalSum(c *fiber.Ctx) error {
	var req CalculateTotalSumRequest

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

	totalSales := h.salesService.GetTotalSum(req.StoreID, startDate, endDate)

	return c.JSON(CalculateTotalSumResponse{
		StoreID:    req.StoreID,
		TotalSales: totalSales,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	})
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
