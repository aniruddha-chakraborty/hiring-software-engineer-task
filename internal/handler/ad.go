package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"sweng-task/internal/model"
	"sweng-task/internal/service"
)

type AdHandler struct {
	logs *zap.SugaredLogger
	ad   *service.AdService
}

var validate = validator.New()

func NewAdHandler(log *zap.SugaredLogger, adv *service.AdService) *AdHandler {
	return &AdHandler{
		logs: log,
		ad:   adv,
	}
}

func (a *AdHandler) GetWinningAds(c *fiber.Ctx) error {
	var query model.WinningAdsQuery

	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Failed to parse query parameters",
			"details": err.Error(),
		})
	}

	if err := validate.Struct(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid query parameters",
			"details": err.Error(),
		})
	}

	advertisements, err := a.ad.GetAd(query.Placement, query.Keyword, query.Category, query.Limit)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "No matching ads found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(advertisements)
}
