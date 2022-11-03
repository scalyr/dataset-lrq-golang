package lrq

import (
	"context"
	"encoding/json"
	"time"
)

type LogRequestAttribs struct {
	StartTime *time.Time
	EndTime   *time.Time
	Filter    *string
	Limit     *uint
}

type logRequest struct {
	QueryType string         `json:"queryType"`
	StartTime *int64         `json:"startTime,omitempty"`
	EndTime   *int64         `json:"endTime,omitempty"`
	Log       logRequestOpts `json:"log"`
}

type logRequestOpts struct {
	Filter *string `json:"filter,omitempty"`
	Limit  *uint   `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}

type logResponseMatch struct {
	LogResponseMatch
	Cursor string `json:"cursor"`
}

type LogResponseMatch struct {
	ServerInfo map[string]interface{} `json:"serverInfo"`
	SessionId  string                 `json:"sessionId"`
	Severity   int                    `json:"severity"`
	ThreadId   string                 `json:"threadId"`
	Timestamp  int64                  `json:"timestamp"`
	Values     map[string]interface{} `json:"values"`
}

func (c *Client) DoLogRequest(ctx context.Context, attribs LogRequestAttribs) ([]LogResponseMatch, error) {
	reqBody := logRequest{
		QueryType: "LOG",
		StartTime: timeToInt64(attribs.StartTime),
		EndTime:   timeToInt64(attribs.EndTime),
		Log: logRequestOpts{
			Filter: attribs.Filter,
			Limit:  attribs.Limit,
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
		Matches []LogResponseMatch `json:"matches"`
	}

	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, err
	}

	return respBody.Matches, nil
}

func (c *Client) DoLogRequestPaginated(ctx context.Context, attribs LogRequestAttribs, cursor *string) ([]LogResponseMatch, *string, error) {
	reqBody := logRequest{
		QueryType: "LOG",
		StartTime: timeToInt64(attribs.StartTime),
		EndTime:   timeToInt64(attribs.EndTime),
		Log: logRequestOpts{
			Filter: attribs.Filter,
			Limit:  attribs.Limit,
			Cursor: cursor,
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	respBytes, err := c.doRequest(ctx, reqBytes)
	if err != nil {
		return nil, nil, err
	}

	var respBody struct {
		Matches []logResponseMatch `json:"matches"`
	}

	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, nil, err
	}

	matches := make([]LogResponseMatch, len(respBody.Matches))
	for i := range respBody.Matches {
		matches[i] = respBody.Matches[i].LogResponseMatch
	}

	return matches, &respBody.Matches[0].Cursor, nil
}
