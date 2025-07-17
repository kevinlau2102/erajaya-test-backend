package dto

type (
	PaginationRequest struct {
		Search string `form:"search"`
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		SortBy string `form:"sortBy"`
		Order  string `form:"order"`
	}

	PaginationResponse struct {
		Page    int   `json:"page"`
		Limit   int   `json:"limit"`
		MaxPage int64 `json:"max_page"`
		Count   int64 `json:"count"`
	}
)

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *PaginationRequest) GetLimit() int {
	return p.Limit
}

func (p *PaginationRequest) GetPage() int {
	return p.Page
}

func (p *PaginationRequest) Default() {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.Limit == 0 {
		p.Limit = 10
	}
}
