package lrq

import (
	"context"
	"encoding/json"
	"time"
)

type TopFacetsRequestAttribs struct {
	StartTime         *time.Time
	EndTime           *time.Time
	Filter            *string
	NumValuesPerFacet *uint
	NumFacets         *uint
}

type topFacetsRequest struct {
	QueryType string               `json:"queryType"`
	StartTime *int64               `json:"startTime,omitempty"`
	EndTime   *int64               `json:"endTime,omitempty"`
	TopFacets topFacetsRequestOpts `json:"topFacets"`
}

type topFacetsRequestOpts struct {
	Filter            *string `json:"filter,omitempty"`
	NumValuesPerFacet *uint   `json:"numValuesToReturnPerFacet,omitempty"`
	NumFacets         *uint   `json:"numFacetsToReturn,omitempty"`
}

type TopFacetsResponseFacet struct {
	Name        string                        `json:"name"`
	MatchCount  float64                       `json:"matchCount"`
	UniqueCount uint                          `json:"uniqueValuesCount"`
	Values      []TopFacetsResponseFacetValue `json:"values"`
}

type TopFacetsResponseFacetValue struct {
	Count uint        `json:"count"`
	Value interface{} `json:"value"`
}

func (c *Client) DoTopFacetsRequest(ctx context.Context, attribs TopFacetsRequestAttribs) ([]TopFacetsResponseFacet, error) {
	reqBody := topFacetsRequest{
		QueryType: "TOP_FACETS",
		StartTime: timeToInt64(attribs.StartTime),
		EndTime:   timeToInt64(attribs.EndTime),
		TopFacets: topFacetsRequestOpts{
			Filter:            attribs.Filter,
			NumValuesPerFacet: attribs.NumValuesPerFacet,
			NumFacets:         attribs.NumFacets,
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.doRequest(ctx, reqBytes)
	if err != nil {
		return nil, err
	}

	var respBody struct {
		Facets []TopFacetsResponseFacet `json:"facets"`
	}

	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, err
	}

	return respBody.Facets, nil
}

type FacetValuesRequestAttribs struct {
	StartTime *time.Time
	EndTime   *time.Time
	Filter    *string
	MaxValues *uint
}

type facetValuesRequest struct {
	QueryType   string                 `json:"queryType"`
	StartTime   *int64                 `json:"startTime,omitempty"`
	EndTime     *int64                 `json:"endTime,omitempty"`
	FacetValues facetValuesRequestOpts `json:"facetValues"`
}

type facetValuesRequestOpts struct {
	Name      string  `json:"name"`
	Filter    *string `json:"filter,omitempty"`
	MaxValues *uint   `json:"maxValues,omitempty"`
}

type FacetValuesResponseValue struct {
	Count uint        `json:"count"`
	Value interface{} `json:"value"`
}

func (c *Client) DoFacetValuesRequest(ctx context.Context, name string, attribs FacetValuesRequestAttribs) ([]FacetValuesResponseValue, error) {
	reqBody := facetValuesRequest{
		QueryType:   "FACET_VALUES",
		StartTime:   timeToInt64(attribs.StartTime),
		EndTime:     timeToInt64(attribs.EndTime),
		FacetValues: facetValuesRequestOpts{
			Name:      name,
			Filter:    attribs.Filter,
			MaxValues: attribs.MaxValues,
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.doRequest(ctx, reqBytes)
	if err != nil {
		return nil, err
	}

	var respBody struct {
		Facet struct {
			MatchCount  float64                    `json:"matchCount"`
			UniqueCount uint                       `json:"uniqueValuesCount"`
			Values      []FacetValuesResponseValue `json:"values"`
		} `json:"facet"`
	}

	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, err
	}

	return respBody.Facet.Values, nil
}
