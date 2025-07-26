package service

import (
	"go.uber.org/zap"
	"sweng-task/internal/model"
)

type DataGeneratorService struct {
	log *zap.SugaredLogger
	lis *LineItemService
}

func NewDataGenerator(log *zap.SugaredLogger, lis *LineItemService) *DataGeneratorService {
	return &DataGeneratorService{
		log: log,
		lis: lis,
	}
}

func (d *DataGeneratorService) GenerateLineItems() {
	// Calls to the helper function are now active.
	d.generateSingleLineItem("Summer Sale Banner", "adv123", 2.5, 3000.0, "homepage_top", []string{"electronics", "sports"}, []string{"summer", "discount"})
	d.generateSingleLineItem("Winter Clearance Promo", "adv456", 3.0, 7000.0, "video_preroll", []string{"fashion", "beauty"}, []string{"clearance", "deal"})
	d.generateSingleLineItem("Travel Deals Campaign", "adv789", 1.8, 5000.0, "article_inline_1", []string{"travel", "food"}, []string{"exclusive", "trending"})
	d.generateSingleLineItem("Gaming Weekend Blast", "adv321", 0.2, 8000.0, "homepage_top", []string{"gaming", "electronics"}, []string{"sale", "new"})
	d.generateSingleLineItem("Home Essentials Discount", "adv654", 2.2, 10000.0, "video_preroll", []string{"home", "sports"}, []string{"deal", "discount"})
	d.generateSingleLineItem("Back to School Deals", "adv111", 3.5, 4000.0, "article_inline_1", []string{"electronics", "fashion"}, []string{"exclusive", "sale"})
	d.generateSingleLineItem("Spring Fashion Promo", "adv222", 1.9, 9000.0, "homepage_top", []string{"fashion", "beauty"}, []string{"trending", "new"})
	d.generateSingleLineItem("Holiday Travel Specials", "adv333", 2.7, 6000.0, "video_preroll", []string{"travel", "food"}, []string{"deal", "exclusive"})
	d.generateSingleLineItem("Fitness Gear Discount", "adv444", 2.3, 3000.0, "article_inline_1", []string{"sports", "home"}, []string{"discount", "sale"})
	d.generateSingleLineItem("Luxury Beauty Sale", "adv555", 4.5, 7000.0, "homepage_top", []string{"beauty", "fashion"}, []string{"clearance", "exclusive"})
	d.generateSingleLineItem("Gadget Madness", "adv666", 3.1, 5000.0, "video_preroll", []string{"electronics", "gaming"}, []string{"trending", "new"})
	d.generateSingleLineItem("Healthy Living Promo", "adv777", 2.6, 8000.0, "article_inline_1", []string{"food", "home"}, []string{"deal", "sale"})
	d.generateSingleLineItem("Weekend Getaway Deals", "adv888", 1.7, 9000.0, "homepage_top", []string{"travel", "sports"}, []string{"exclusive", "discount"})
	d.generateSingleLineItem("Clearance Electronics", "adv999", 3.8, 1000.0, "video_preroll", []string{"electronics", "home"}, []string{"clearance", "deal"})
	d.generateSingleLineItem("Gaming Console Offer", "adv112", 4.0, 6000.0, "article_inline_1", []string{"gaming", "electronics"}, []string{"new", "sale"})
	d.generateSingleLineItem("Cozy Home Sale", "adv113", 2.4, 2000.0, "homepage_top", []string{"home", "fashion"}, []string{"discount", "exclusive"})
	d.generateSingleLineItem("Fashion Week Promo", "adv114", 2.9, 3000.0, "video_preroll", []string{"fashion", "beauty"}, []string{"trending", "sale"})
	d.generateSingleLineItem("Sports Gear Blowout", "adv115", 3.3, 5000.0, "article_inline_1", []string{"sports", "gaming"}, []string{"deal", "clearance"})
	d.generateSingleLineItem("Smart Home Specials", "adv116", 2.1, 7000.0, "homepage_top", []string{"electronics", "home"}, []string{"exclusive", "new"})
	d.generateSingleLineItem("Travel Light Promo", "adv117", 1.5, 4000.0, "video_preroll", []string{"travel", "fashion"}, []string{"discount", "sale"})
}

func (d *DataGeneratorService) generateSingleLineItem(name string, advID string, bid, budget float64, placement string, categories []string, keywords []string) {
	input := model.LineItemCreate{
		Name:         name,
		AdvertiserID: advID,
		Bid:          bid,
		Budget:       budget,
		Placement:    placement,
		Categories:   categories,
		Keywords:     keywords,
	}
	_, err := d.lis.Create(input)
	if err != nil {
		d.log.Error("Failed to create lineItem", zap.Error(err))
		return
	}
}
