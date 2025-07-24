package model

type WinningAdsQuery struct {
	Placement string `query:"placement" validate:"required,min=3,max=50"`
	Keyword   string `query:"keyword" validate:"omitempty,max=50"`
	Category  string `query:"category" validate:"omitempty,max=50"`
}
