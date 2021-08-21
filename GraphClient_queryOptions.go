package msgraph

import (
	"context"
	"net/http"
	"net/url"
)

type getRequestParams interface {
	Context() context.Context
	Values() url.Values
	Headers() http.Header
}

type GetQueryOption func(opts *getQueryOptions)

type ListQueryOption func(opts *listQueryOptions)

type CreateQueryOption func(opts *createQueryOptions)

type UpdateQueryOption func(opts *updateQueryOptions)

type DeleteQueryOption func(opts *deleteQueryOptions)

var (
	// GetWithContext - add a context.Context to the HTTP request e.g. to allow cancellation
	GetWithContext = func(ctx context.Context) GetQueryOption {
		return func(opts *getQueryOptions) {
			opts.ctx = ctx
		}
	}

	// GetWithSelect - $select - Filters properties (columns) - https://docs.microsoft.com/en-us/graph/query-parameters#select-parameter
	GetWithSelect = func(selectParam string) GetQueryOption {
		return func(opts *getQueryOptions) {
			opts.queryValues.Add(odataSelectParamKey, selectParam)
		}
	}

	// ListWithContext - add a context.Context to the HTTP request e.g. to allow cancellation
	ListWithContext = func(ctx context.Context) ListQueryOption {
		return func(opts *listQueryOptions) {
			opts.ctx = ctx
		}
	}

	// ListWithSelect - $select - Filters properties (columns) - https://docs.microsoft.com/en-us/graph/query-parameters#select-parameter
	ListWithSelect = func(selectParam string) ListQueryOption {
		return func(opts *listQueryOptions) {
			opts.queryValues.Add(odataSelectParamKey, selectParam)
		}
	}

	// ListWithFilter - $filter - Filters results (rows) - https://docs.microsoft.com/en-us/graph/query-parameters#filter-parameter
	ListWithFilter = func(filterParam string) ListQueryOption {
		return func(opts *listQueryOptions) {
			opts.queryValues.Add(odataFilterParamKey, filterParam)
		}
	}

	// ListWithSearch - $search - Returns results based on search criteria - https://docs.microsoft.com/en-us/graph/query-parameters#search-parameter
	ListWithSearch = func(searchParam string) ListQueryOption {
		return func(opts *listQueryOptions) {
			opts.queryHeaders.Add("ConsistencyLevel", "eventual")
			opts.queryValues.Add(odataSearchParamKey, searchParam)
		}
	}

	// CreateWithContext - add a context.Context to the HTTP request e.g. to allow cancellation
	CreateWithContext = func(ctx context.Context) CreateQueryOption {
		return func(opts *createQueryOptions) {
			opts.ctx = ctx
		}
	}

	// UpdateWithContext - add a context.Context to the HTTP request e.g. to allow cancellation
	UpdateWithContext = func(ctx context.Context) UpdateQueryOption {
		return func(opts *updateQueryOptions) {
			opts.ctx = ctx
		}
	}
	// DeleteWithContext - add a context.Context to the HTTP request e.g. to allow cancellation
	DeleteWithContext = func(ctx context.Context) DeleteQueryOption {
		return func(opts *deleteQueryOptions) {
			opts.ctx = ctx
		}
	}
)

// getQueryOptions allow to optionally pass OData query options
// see https://docs.microsoft.com/en-us/graph/query-parameters
type getQueryOptions struct {
	ctx         context.Context
	queryValues url.Values
}

func (g *getQueryOptions) Context() context.Context {
	if g.ctx == nil {
		return context.Background()
	}
	return g.ctx
}

func (g getQueryOptions) Values() url.Values {
	return g.queryValues
}

func (g getQueryOptions) Headers() http.Header {
	return http.Header{}
}

func compileGetQueryOptions(options []GetQueryOption) *getQueryOptions {
	var opts = &getQueryOptions{
		queryValues: url.Values{},
	}
	for idx := range options {
		options[idx](opts)
	}

	return opts
}

// listQueryOptions allow to optionally pass OData query options
// see https://docs.microsoft.com/en-us/graph/query-parameters
type listQueryOptions struct {
	getQueryOptions
	queryHeaders http.Header
}

func (g *listQueryOptions) Context() context.Context {
	if g.ctx == nil {
		return context.Background()
	}
	return g.ctx
}

func (g listQueryOptions) Values() url.Values {
	return g.queryValues
}

func (g listQueryOptions) Headers() http.Header {
	return g.queryHeaders
}

func compileListQueryOptions(options []ListQueryOption) *listQueryOptions {
	var opts = &listQueryOptions{
		getQueryOptions: getQueryOptions{
			queryValues: url.Values{},
		},
		queryHeaders: http.Header{},
	}
	for idx := range options {
		options[idx](opts)
	}

	return opts
}

// createQueryOptions allows to add a context to the request
type createQueryOptions struct {
	getQueryOptions
}

func (g *createQueryOptions) Context() context.Context {
	if g.ctx == nil {
		return context.Background()
	}
	return g.ctx
}

func compileCreateQueryOptions(options []CreateQueryOption) *createQueryOptions {
	var opts = &createQueryOptions{
		getQueryOptions: getQueryOptions{
			queryValues: url.Values{},
		},
	}
	for idx := range options {
		options[idx](opts)
	}

	return opts
}

// updateQueryOptions allows to add a context to the request
type updateQueryOptions struct {
	getQueryOptions
}

func (g *updateQueryOptions) Context() context.Context {
	if g.ctx == nil {
		return context.Background()
	}
	return g.ctx
}

func compileUpdateQueryOptions(options []UpdateQueryOption) *updateQueryOptions {
	var opts = &updateQueryOptions{
		getQueryOptions: getQueryOptions{
			queryValues: url.Values{},
		},
	}
	for idx := range options {
		options[idx](opts)
	}

	return opts
}

// deleteQueryOptions allows to add a context to the request
type deleteQueryOptions struct {
	getQueryOptions
}

func (g *deleteQueryOptions) Context() context.Context {
	if g.ctx == nil {
		return context.Background()
	}
	return g.ctx
}

func compileDeleteQueryOptions(options []DeleteQueryOption) *deleteQueryOptions {
	var opts = &deleteQueryOptions{
		getQueryOptions: getQueryOptions{
			queryValues: url.Values{},
		},
	}
	for idx := range options {
		options[idx](opts)
	}

	return opts
}
