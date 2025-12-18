package dtos

type PaginatedUsersQueryDto struct {
	Page           int    `query:"page" validate:"omitempty,min=1"`
	Limit          int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search         string `query:"search" validate:"omitempty,max=255"`
	SearchField    string `query:"searchField" validate:"omitempty,oneof=username email"`
	Order          string `query:"order" validate:"omitempty,oneof=asc desc"`
	SortBy         string `query:"sortBy" validate:"omitempty,oneof=created_at updated_at username"`
	IncludeDeleted bool   `query:"includeDeleted"`
}

func (q *PaginatedUsersQueryDto) ApplyDefaults() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 5
	}
	if q.SearchField == "" {
		q.SearchField = "username"
	}
	if q.Order == "" {
		q.Order = "desc"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
}
