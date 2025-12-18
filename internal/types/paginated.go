package types

type Paginated[T any] struct {
	Data        []T
	TotalPages  int
	CurrentPage int
	Limit       int
	Order       string
	SortBy      string
	HasPrevious bool
	HasNext     bool
}
