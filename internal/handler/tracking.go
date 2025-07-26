package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"sweng-task/internal/model"
	"sweng-task/internal/service"
	"time"
)

type TrackingHandler struct {
	logs   *zap.SugaredLogger
	pubSub *service.PubSub
}

func NewTrackingHandler(log *zap.SugaredLogger, sub *service.PubSub) *TrackingHandler {
	return &TrackingHandler{
		logs:   log,
		pubSub: sub,
	}
}

func (t *TrackingHandler) TrackEvent(c *fiber.Ctx) error {
	var query model.TrackingEvent
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Failed to parse query parameters",
			"details": err.Error(),
		})
	}
	if err := c.BodyParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid JSON"})
	}
	if err := validate.Struct(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid query parameters",
			"details": err.Error(),
		})
	}

	stats := map[string]int{}
	if query.EventType == "impression" {
		stats["impression"]++
	} else if query.EventType == "click" {
		stats["click"]++
	} else if query.EventType == "conversion" {
		stats["conversion"]++
	}

	msg := fmt.Sprintf(`{"event_time": %q,"event_minute": %d,"item_id": "%s","user_id": "%s","placement": "%s","keyword": "%s","clicks": %d,"impressions": %d,"conversions": %d}`, time.Now().Format(time.RFC3339Nano), time.Now().Minute(), query.LineItemID, query.UserID,
		query.Placement, "keyword", stats["impression"], stats["click"], stats["conversion"])

	// conversion to string possible from query model.TrackingEvent to json, but then i have to sarama have change it to byte
	// that's same thing, that's why msg was built like this.
	t.pubSub.Publish(msg)
	return c.JSON(fiber.StatusAccepted)
}
