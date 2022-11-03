package lrq

import (
	"context"
	"encoding/json"
	"time"
)

type PQRequestAttribs struct {
	StartTime *time.Time
	EndTime   *time.Time
}

type pqRequest struct {
	QueryType string        `json:"queryType"`
	StartTime *int64        `json:"startTime,omitempty"`
	EndTime   *int64        `json:"endTime,omitempty"`
	PQ        pqRequestOpts `json:"pq"`
}

type pqRequestOpts struct {
	Query      string `json:"query"`
	ResultType string `json:"resultType"`
}

func (c *Client) DoPQPlotRequest(ctx context.Context, query string, attribs PQRequestAttribs) (*PlotResponseData, error) {
	reqBody := pqRequest{
		QueryType: "PQ",
		StartTime: timeToInt64(attribs.StartTime),
		EndTime:   timeToInt64(attribs.EndTime),
		PQ: pqRequestOpts{
			Query:      query,
			ResultType: "PLOT",
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

type PQResponseData struct {
	Values  [][]interface{} `json:"values"`
	Columns []struct {
		Name string `json:"name"`
		Type string `json:"cellType"`
	} `json:"columns"`
	Warnings []string `json:"warnings"`
}

func (c *Client) DoPQTableRequest(ctx context.Context, query string, attribs PQRequestAttribs) (*PQResponseData, error) {
	reqBody := pqRequest{
		QueryType: "PQ",
		StartTime: timeToInt64(attribs.StartTime),
		EndTime:   timeToInt64(attribs.EndTime),
		PQ: pqRequestOpts{
			Query:      query,
			ResultType: "TABLE",
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

	var respBody PQResponseData
	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}
