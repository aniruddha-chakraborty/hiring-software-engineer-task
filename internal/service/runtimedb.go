package service

import "go.uber.org/zap"

type RunTimeDB struct {
	log            *zap.SugaredLogger
	Keywords       map[string][]string
	Categories     map[string][]string
	Placements     map[string][]string
	TargetFree     map[string]float32
	ParameterCount map[string]int
}

func NewRunTimeDB(log *zap.SugaredLogger) *RunTimeDB {
	return &RunTimeDB{
		log:            log,
		Keywords:       map[string][]string{},
		Categories:     map[string][]string{},
		Placements:     map[string][]string{},
		TargetFree:     map[string]float32{},
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

func (r *RunTimeDB) AddCategory(Categories []string, advertisementId string) {
	for _, category := range Categories {
		if _, ok := r.Categories[category]; ok {
			r.Categories[category] = append(r.Categories[category], advertisementId)
		} else {
			r.Categories[category] = []string{advertisementId}
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

func (r *RunTimeDB) AddPlacements(placement string, advertisementId string) {
	if _, ok := r.Placements[placement]; ok {
		r.Placements[placement] = append(r.Placements[placement], advertisementId)
	} else {
		r.Placements[placement] = []string{advertisementId}
	}
}

func (r *RunTimeDB) GetPlacements(placement string) []string {
	if _, ok := r.Placements[placement]; ok {
		return r.Placements[placement]
	} else {
		return []string{}
	}
}

func (r *RunTimeDB) AddTargetFree(advertisementId string) {
	if _, ok := r.TargetFree[advertisementId]; ok {
		r.TargetFree[advertisementId] = 0.0
	}
}

func (r *RunTimeDB) GetInitialScoringWithTargetFreeItems() map[string]float32 {
	return r.TargetFree
}

func (r *RunTimeDB) AddParameterCount(advertiserId string, parameters int) {
	r.ParameterCount[advertiserId] = parameters
}

func (r *RunTimeDB) GetParameterCount() map[string]int {
	return r.ParameterCount
}
