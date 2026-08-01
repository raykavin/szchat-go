package szchat

import (
	"net/url"
	"strconv"
)

// PaginatedResponse is the Laravel-style paginator envelope returned by most
// SZChat list endpoints.
type PaginatedResponse[T any] struct {
	CurrentPage  int     `json:"current_page"`
	Data         []T     `json:"data"`
	Total        int     `json:"total"`
	PerPage      int     `json:"per_page"`
	LastPage     int     `json:"last_page,omitempty"`
	From         int     `json:"from,omitempty"`
	To           int     `json:"to,omitempty"`
	Path         string  `json:"path,omitempty"`
	FirstPageURL string  `json:"first_page_url,omitempty"`
	LastPageURL  string  `json:"last_page_url,omitempty"`
	NextPageURL  *string `json:"next_page_url,omitempty"`
	PrevPageURL  *string `json:"prev_page_url,omitempty"`
}

// ListOptions holds the pagination query parameters accepted by most list
// endpoints.
type ListOptions struct {
	Page     int
	Limit    int
	Paginate string
}

func (o ListOptions) values() url.Values {
	q := url.Values{}
	setIntParam(q, "page", o.Page)
	setIntParam(q, "limit", o.Limit)
	setParam(q, "paginate", o.Paginate)
	return q
}

func setParam(q url.Values, key, value string) {
	if value != "" {
		q.Set(key, value)
	}
}

func setIntParam(q url.Values, key string, value int) {
	if value > 0 {
		q.Set(key, strconv.Itoa(value))
	}
}
