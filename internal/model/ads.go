package model

// Ad represents an advertisement ready to be served
type Ad struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	AdvertiserID string  `json:"advertiser_id"`
	Bid          float64 `json:"bid"`
	Placement    string  `json:"placement"`
	ServeURL     string  `json:"serve_url"`
	Relevance    int     `json:"relevance"`
}

// WinningAdsQuery represents Winning ad request from router and specifies its requirement
type WinningAdsQuery struct {
	Placement string `query:"placement" validate:"required,max=50"`
	Keyword   string `query:"keyword" validate:"omitempty,max=50"`
	Category  string `query:"category" validate:"omitempty,max=50"`
}
