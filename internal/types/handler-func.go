package types

import "net/http"

type HandlerFunc func(r *http.Request) (any, error)
