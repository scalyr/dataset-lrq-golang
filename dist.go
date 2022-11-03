package lrq

import (
	"context"
	"encoding/json"
	"time"
)

type DistRequestAttribs struct {
	StartTime *time.Time
	EndTime   *time.Time
	Filter    *string
	Facet     *string
}

type distRequest struct {
	QueryType string          `json:"queryType"`
	StartTime *int64          `json:"startTime,omitempty"`
	EndTime   *int64          `json:"endTime,omitempty"`
	Dist      distRequestOpts `json:"distribution"`
}

type distRequestOpts struct {
	Filter *string `json:"filter,omitempty"`
	Facet  *string `json:"facet,omitempty"`
}

type DistResponseData struct {
	PositiveSampleCounts []float64 `json:"positiveSampleCounts"`
	NegativeSampleCounts []float64 `json:"negativeSampleCounts"`
	PositiveXAxis        []float64 `json:"positiveXAxis"`
	NegativeXAxis        []float64 `json:"negativeXAxis"`
	PositiveBucketWidths []float64 `json:"positiveBucketWidths"`
	NegativeBucketWidths []float64 `json:"negativeBucketWidths"`

	ZeroCount            float64   `json:"zeroCount"`
	Mean                 float64   `json:"mean"`
	Median               float64   `json:"median"`
	Min                  float64   `json:"min"`
	Max                  float64   `json:"max"`
	P10                  float64   `json:"p10"`
	P90                  float64   `json:"p90"`
	P99                  float64   `json:"p99"`
	P999                 float64   `json:"p999"`
}

func (c *Client) DoDistRequest(ctx context.Context, attribs DistRequestAttribs) (*DistResponseData, error) {
	reqBody := distRequest{
		QueryType: "DISTRIBUTION",
		StartTime: timeToInt64(attribs.StartTime),
		EndTime:   timeToInt64(attribs.EndTime),
		Dist:      distRequestOpts{
			Filter: attribs.Filter,
			Facet:  attribs.Facet,
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

	var respBody DistResponseData
	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}
