package dtos

type PaginatedProductsQueryDto struct {
	Page           int    `query:"page" validate:"omitempty,min=1"`
	Limit          int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search         string `query:"search" validate:"omitempty,max=255"`
	SearchField    string `query:"searchField" validate:"omitempty,oneof=name"`
	Order          string `query:"order" validate:"omitempty,oneof=asc desc"`
	SortBy         string `query:"sortBy" validate:"omitempty,oneof=created_at updated_at name"`
	IncludeDeleted bool   `query:"includeDeleted"`
}

func (q *PaginatedProductsQueryDto) ApplyDefaults() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 5
	}
	if q.SearchField == "" {
		q.SearchField = "name"
	}
	if q.Order == "" {
		q.Order = "desc"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
}
