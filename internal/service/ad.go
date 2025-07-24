package service

import (
	"fmt"
	"go.uber.org/zap"
	"math"
	"sweng-task/internal/model"
)

var CoreScoring map[string]float32 = map[string]float32{
	"keywordWeight":  2,
	"categoryWeight": 5,
	"bidWeight":      10,
	"paramWeight":    5,
}

type AdService struct {
	logs      *zap.SugaredLogger
	runTimeDB *RunTimeDB
	lis       *LineItemService
}

func NewAdService(log *zap.SugaredLogger, runTimeDB *RunTimeDB, lis *LineItemService) *AdService {
	return &AdService{
		logs:      log,
		runTimeDB: runTimeDB,
		lis:       lis,
	}
}

// This whole thing optimises the FindMatchingLineItems and the ad selection part together, It's much more efficient
func (s *AdService) GetAd(placement string, keyword string, category string, limit int) ([]*model.Ad, error) {
	if len(s.runTimeDB.GetPlacements(placement)) == 0 {
		return []*model.Ad{}, nil
	}

	highestBid := -1.0
	highestBidderId := "DummyString"
	// Applying bucket sort, because as i need limited number of ads, so at scale (5K line items) we don't need to sort the whole array
	// we just need to sort the highest value bucket until we hit desired amount of result. plus bucket sort is O(N) in insert and retrival
	// does not need O(N) at all! We can teak the minBin,maxBid and bucketGap later to increase the performance
	minBid := 0.01
	maxBid := 10.0
	bucketGap := 0.50
	bucketCount := int(math.Ceil((maxBid - minBid) / bucketGap))
	buckets := make([][]*model.LineItem, bucketCount)

	insertIntoBucket := func(item *model.LineItem) {
		idx := int((item.Bid - minBid) / bucketGap)
		if idx < 0 {
			idx = 0
		}
		if idx >= bucketCount {
			idx = bucketCount - 1
		}
		buckets[idx] = append(buckets[idx], item)
	}
	// bucket sort prep ends
	score := s.runTimeDB.GetInitialScoringWithTargetFreeItems()
	paramMatch := map[string]int{}
	// Initial scoring loop
	for id := range score {
		highestBid, highestBidderId = s.updateHighestBid(highestBid, highestBidderId, id)
		insertIntoBucket(s.lis.items[id])
	}

	// Keyword scoring loop
	keywordIds := s.runTimeDB.GetKeyWords(keyword)
	for _, id := range keywordIds {
		score[id] += CoreScoring["keywordWeight"]
		highestBid, highestBidderId = s.updateHighestBid(highestBid, highestBidderId, id)
		// I need to check what percent of a particular line item is getting matched, so that i can send back in relvence
		paramMatch[id]++
		// if all params match that means the ad the 100% relevant
		if paramMatch[id] == s.runTimeDB.ParameterCount[id] {
			score[id] += CoreScoring["paramWeight"]
		}
		insertIntoBucket(s.lis.items[id])
	}

	// Category scoring loop
	categoryIds := s.runTimeDB.GetCategory(category)
	for _, id := range categoryIds {
		score[id] += CoreScoring["categoryWeight"]
		highestBid, highestBidderId = s.updateHighestBid(highestBid, highestBidderId, id)
		paramMatch[id]++
		// if all params match that means the ad the 100% relevant
		if paramMatch[id] == s.runTimeDB.ParameterCount[id] {
			score[id] += CoreScoring["paramWeight"]
		}
		insertIntoBucket(s.lis.items[id])
	}

	if highestBidderId != "DummyString" {
		// priority scoring happened
		score[highestBidderId] += CoreScoring["bidWeight"]
	}

	result := []*model.Ad{}
	for i := bucketCount - 1; i >= 0 && len(result) < limit; i-- {
		for _, ad := range buckets[i] {
			result = append(result, &model.Ad{
				ID:           ad.ID,
				Name:         s.lis.items[ad.ID].Name,
				AdvertiserID: s.lis.items[ad.ID].AdvertiserID,
				Bid:          s.lis.items[ad.ID].Bid,
				Placement:    s.lis.items[ad.ID].Placement,
				ServeURL:     fmt.Sprintf("https://content.realtimemediatool.com/data/%s", ad.ID),
				Relevance:    (paramMatch[ad.ID] * 100) / s.runTimeDB.ParameterCount[ad.ID],
			})
			if len(result) == limit {
				return result, nil
			}
		}
	}
	return result, nil
}

func (s *AdService) updateHighestBid(currentHighest float64, currentBidder string, candidateID string) (float64, string) {
	bid := s.lis.items[candidateID].Bid
	if bid > currentHighest {
		return bid, candidateID
	}
	return currentHighest, currentBidder
}
