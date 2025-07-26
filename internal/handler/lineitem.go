package handler

import (
	"sweng-task/internal/model"

	"sweng-task/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// LineItemHandler handles HTTP requests related to line items
type LineItemHandler struct {
	service *service.LineItemService
	log     *zap.SugaredLogger
}

// NewLineItemHandler creates a new LineItemHandler
func NewLineItemHandler(service *service.LineItemService, log *zap.SugaredLogger) *LineItemHandler {
	return &LineItemHandler{
		service: service,
		log:     log,
	}
}

// Create handles the creation of a new line item
func (h *LineItemHandler) Create(c *fiber.Ctx) error {
	var input model.LineItemCreate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Note: Validation logic should be implemented by the candidate
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid query parameters",
			"details": err.Error(),
		})
	}
	lineItem, err := h.service.Create(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to create line item",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(lineItem)
}

// GetByID handles retrieving a line item by ID
func (h *LineItemHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Missing line item ID",
		})
	}

	lineItem, err := h.service.GetByID(id)
	if err != nil {
		if err == service.ErrLineItemNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    fiber.StatusNotFound,
				"message": "Line item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to retrieve line item",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(lineItem)
}

// GetAll handles retrieving all line items with optional filtering
func (h *LineItemHandler) GetAll(c *fiber.Ctx) error {
	advertiserID := c.Query("advertiser_id")
	placement := c.Query("placement")

	lineItems, err := h.service.GetAll(advertiserID, placement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to retrieve line items",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(lineItems)
}
