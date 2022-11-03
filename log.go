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
	Filter    *string `json:"filter,omitempty"`
	Limit     *uint   `json:"limit,omitempty"`
	Cursor    *string `json:"cursor,omitempty"`
	Ascending *bool   `json:"ascending,omitempty"`
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

			// This is needed for cursor matching later
			Ascending: func() *bool { b := true; return &b }(),
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

	// The cursor represents an identifier associated with a specific event.
	// So specifying a cursor means the matches should start with that event,
	// which here is the last event of the previous paginated results.

	if cursor == nil {
		matches := make([]LogResponseMatch, len(respBody.Matches))
		for i := range respBody.Matches {
			matches[i] = respBody.Matches[i].LogResponseMatch
		}

		lastCursor := respBody.Matches[len(respBody.Matches)-1].Cursor
		return matches, &lastCursor, nil
	} else {
		if len(respBody.Matches) == 0 {
			return []LogResponseMatch{}, nil, nil
		}

		var matches []LogResponseMatch
		firstCursor := respBody.Matches[0].Cursor
		if firstCursor == *cursor {
			matches = make([]LogResponseMatch, len(respBody.Matches)-1)
			for i := range respBody.Matches[1:] {
				matches[i] = respBody.Matches[i+1].LogResponseMatch
			}
		} else {
			matches = make([]LogResponseMatch, len(respBody.Matches))
			for i := range respBody.Matches {
				matches[i] = respBody.Matches[i].LogResponseMatch
			}
		}

		lastCursor := respBody.Matches[len(respBody.Matches)-1].Cursor
		return matches, &lastCursor, nil
	}
}
