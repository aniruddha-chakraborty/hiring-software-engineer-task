package service

import "go.uber.org/zap"

type RunTimeDB struct {
	log            *zap.SugaredLogger
	Keywords       map[string][]string
	Categories     map[string][]string
	Placements     map[string][]string
	TargetFree     map[string]float64
	ParameterCount map[string]int
}

func NewRunTimeDB(log *zap.SugaredLogger) *RunTimeDB {
	return &RunTimeDB{
		log:            log,
		Keywords:       map[string][]string{},
		Categories:     map[string][]string{},
		Placements:     map[string][]string{},
		TargetFree:     map[string]float64{},
		ParameterCount: map[string]int{},
	}
}

func (r *RunTimeDB) AddKeyWords(Keywords []string, advertisementId string) {
	for _, keyword := range Keywords {
		if _, ok := r.Keywords[keyword]; ok {
			r.Keywords[keyword] = append(r.Keywords[keyword], advertisementId)
		} else {
			r.Keywords[keyword] = []string{advertisementId}
		}
	}
}

func (r *RunTimeDB) GetKeyWords(keyword string) []string {
	if _, ok := r.Keywords[keyword]; ok {
		return r.Keywords[keyword]
	} else {
		return []string{}
	}
}

func (r *RunTimeDB) AddCategory(Categories []string, lineItemId string) {
	for _, category := range Categories {
		if _, ok := r.Categories[category]; ok {
			r.Categories[category] = append(r.Categories[category], lineItemId)
		} else {
			r.Categories[category] = []string{lineItemId}
		}
	}
}

func (r *RunTimeDB) GetCategory(category string) []string {
	if _, ok := r.Categories[category]; ok {
		return r.Categories[category]
	} else {
		return []string{}
	}
}

func (r *RunTimeDB) AddPlacements(placement string, lineItemId string) {
	if _, ok := r.Placements[placement]; ok {
		r.Placements[placement] = append(r.Placements[placement], lineItemId)
	} else {
		r.Placements[placement] = []string{lineItemId}
	}
}

func (r *RunTimeDB) GetPlacements(placement string) []string {
	if _, ok := r.Placements[placement]; ok {
		return r.Placements[placement]
	} else {
		return []string{}
	}
}

func (r *RunTimeDB) AddTargetFree(lineItemId string) {
	if _, ok := r.TargetFree[lineItemId]; ok {
		r.TargetFree[lineItemId] = 0.0
	}
}

func (r *RunTimeDB) GetInitialScoringWithTargetFreeItems() map[string]float64 {
	// Create a new map to be the copy
	scoreCopy := make(map[string]float64, len(r.TargetFree))

	// Loop through the original and populate the copy
	for key, value := range r.TargetFree {
		scoreCopy[key] = value
	}

	// Return the safe copy, not the original
	return scoreCopy
}

func (r *RunTimeDB) AddParameterCount(advertiserId string, parameters int) {
	r.ParameterCount[advertiserId] = parameters
}

func (r *RunTimeDB) GetParameterCount() map[string]int {
	return r.ParameterCount
}
