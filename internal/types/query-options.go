package types

type SelectFields struct {
	SelectFields []string
}

type QueryOptions struct {
	Include      string
	SelectFields *SelectFields
}

type OrderOptions struct {
	WithUser    bool
	WithProduct bool
}
