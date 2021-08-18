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
