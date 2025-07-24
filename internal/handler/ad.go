package handler

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"sweng-task/internal/service"
)

type AdHandler struct {
	logs *zap.SugaredLogger
	ad   *service.AdService
}

func NewAdHandler(log *zap.SugaredLogger, adv *service.AdService) *AdHandler {
	return &AdHandler{
		logs: log,
		ad:   adv,
	}
}

func (a *AdHandler) GetWinningAds(c *fiber.Ctx) error {
	placement := c.Query("placement")
	keyword := c.Query("keyword")
	category := c.Query("category")

	advertisements, err := a.ad.GetAd(placement, keyword, category, 4)
	if err != nil {
		spew.Dump(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "Ad not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(advertisements)
}
