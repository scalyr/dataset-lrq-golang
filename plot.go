package lrq

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type PlotRequestAttribs struct {
	StartTime         *time.Time
	EndTime           *time.Time
	Filter            *string
	Slices            *uint
	SliceWidth        *string
	BreakdownFacet    *string
}

type plotRequest struct {
	QueryType string          `json:"queryType"`
	StartTime *int64          `json:"startTime,omitempty"`
	EndTime   *int64          `json:"endTime,omitempty"`
	Plot      plotRequestOpts `json:"plot"`
}

type plotRequestOpts struct {
	Filter            *string `json:"filter,omitempty"`
	Slices            *uint   `json:"slices,omitempty"`
	SliceWidth        *string `json:"sliceWidth,omitempty"`
	Expression        *string `json:"expression"`
	BreakdownFacet    *string `json:"breakdownFacet,omitempty"`
	Frequency         *string `json:"frequency"`
}

type PlotResponseData struct {
	XAxis []int64 `json:"xAxis"`
	Plots []struct {
		Label *string    `json:"label"`
		Values []float64 `json:"samples"`
	} `json:"plots"`
}

func (c *Client) DoPlotRequest(ctx context.Context, expression string, attribs PlotRequestAttribs) (*PlotResponseData, error) {
	// At least and at most one of Slices or SliceWidth must be defined
	if attribs.Slices == nil && attribs.SliceWidth == nil {
		return nil, errors.New("either Slices or SliceWidth must be defined")
	} else if attribs.Slices != nil && attribs.SliceWidth != nil {
		return nil, errors.New("both Slices and SliceWidth cannot be defined")
	}

	reqBody := plotRequest{
		QueryType: "PLOT",
		StartTime: timeToInt64(attribs.StartTime),
		EndTime:   timeToInt64(attribs.EndTime),
		Plot:      plotRequestOpts{
			Filter:         attribs.Filter,
			Slices:         attribs.Slices,
			SliceWidth:     attribs.SliceWidth,
			Expression:     &expression,
			BreakdownFacet: attribs.BreakdownFacet,
			Frequency:      func() *string { s := "LOW"; return &s }(),
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

	var respBody PlotResponseData
	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}
