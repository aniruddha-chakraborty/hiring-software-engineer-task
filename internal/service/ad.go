package service

import (
	"fmt"
	"go.uber.org/zap"
	"math"
	_ "slices"
	"sort"
	"sweng-task/internal/model"
)

var CoreScoring map[string]float64 = map[string]float64{
	"keywordWeight":  5,
	"categoryWeight": 5,
	"bidWeight":      6,
	"paramWeight":    5,
}

// We can implement catagory specific scoring
var CategoryScoring map[string]float64 = map[string]float64{}

// We can implement keywords specific scoring as well and add with existing implemented scoring
var KeyWordsScoring map[string]float64 = map[string]float64{}

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
	relevanceSore := map[string]int{}
	highestBid := -1.0
	highestBidderId := "DummyString"
	// Applying bucket sort, because as i need limited number of ads, so at scale (5K line items) we don't need to sort the whole array
	// we just need to sort the highest value bucket until we hit desired amount of result. plus bucket sort is O(N) in insert and retrival
	// does not need O(N) at all! We can tweak the minBin,maxBid and bucketGap later to increase the performance
	minBid := 0.01
	maxBid := 10.0
	bucketGap := 0.50
	bucketCount := int(math.Ceil((maxBid - minBid) / bucketGap))
	buckets := make([][]*model.LineItem, bucketCount)

	insertIntoBucket := func(item *model.LineItem, score float64) {
		idx := int((score - minBid) / bucketGap)
		if idx < 0 {
			idx = 0
		}
		if idx >= bucketCount {
			idx = bucketCount - 1
		}
		buckets[idx] = append(buckets[idx], item)
	}
	// bucket sort prep ends
	// There can be lineitems that does not have any targeting, created separate step for this use case
	score := s.runTimeDB.GetInitialScoringWithTargetFreeItems()
	paramMatch := map[string]int{}
	// Initial scoring loop
	for id := range score {
		highestBid, highestBidderId = s.updateHighestBid(highestBid, highestBidderId, id)
		insertIntoBucket(s.lis.items[id], score[id])
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
		insertIntoBucket(s.lis.items[id], score[id])
		relevanceSore[id] += 50
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
		insertIntoBucket(s.lis.items[id], score[id])
		relevanceSore[id] += 50
	}

	if highestBidderId != "DummyString" {
		// priority scoring done
		score[highestBidderId] += CoreScoring["bidWeight"]
	}
	// Even though it seems like This is Only worst case O(N^2), and that worst case is impossible to hit
	result := []*model.Ad{}
	for i := bucketCount - 1; i >= 0 && len(result) < limit; i-- {
		// This is running on very few number of items, that why this sort will is extremly efficient, also i think we can do a pre-sort
		// type stuff during the insertion which will reduce sorting time more in big scale
		sort.Slice(buckets[i], func(a, b int) bool {
			idA := buckets[i][a].ID
			idB := buckets[i][b].ID
			return score[idA] > score[idB]
		})
		for _, ad := range buckets[i] {
			result = append(result, &model.Ad{
				ID:           ad.ID,
				Name:         s.lis.items[ad.ID].Name,
				AdvertiserID: s.lis.items[ad.ID].AdvertiserID,
				Bid:          s.lis.items[ad.ID].Bid,
				Placement:    s.lis.items[ad.ID].Placement,
				ServeURL:     fmt.Sprintf("https://content.realtimemediatool.com/data/%s", ad.ID),
				// Normally, there will be multiple keywords and you need to match with
				Relevance: relevanceSore[ad.ID],
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
