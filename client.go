package lrq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	datasetUrl  string
	apikey      string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

func NewClient(datasetUrl string, apikey string) *Client {
	return &Client{
		datasetUrl: datasetUrl,
		apikey:     apikey,

		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},

		// 100 requests per minute with burst/bucket size of 100
		rateLimiter: rate.NewLimiter(100*rate.Every(1*time.Minute), 100),
	}
}

func (c *Client) makeRequest(method, url string, body io.Reader) (*http.Request, error) {
	// TODO Use runtime/debug.ReadBuildInfo() instead
	const VERSION = "1.0.0"

	if err := c.rateLimiter.Wait(context.Background()); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apikey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dataset-lrq-golang/"+VERSION)

	return req, nil
}

func (c *Client) doRequest(ctx context.Context, body []byte) ([]byte, error) {
	// Long-Running Query (LRQ) api usage:
	// - An initial POST request is made containing the standard/power query
	// - Its response may or may not contain the results
	//   - This is indicated by stepsCompleted == totalSteps in the response
	// - If not complete, follow up with GET ping requests with the response id
	//   - Include the token from the initial POST request response
	// - When complete send a DELETE request to clean up resources
	//   - Include the token from the initial POST request response

	const TOKEN_HEADER = "X-Dataset-Query-Forward-Tag"

	isSuccessful := func(r *http.Response) bool {
		return 200 <= r.StatusCode && r.StatusCode < 300
	}

	req, err := c.makeRequest("POST", c.datasetUrl+"/v2/api/queries", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var respBody struct {
		Id             string          `json:"id"`
		StepsCompleted int             `json:"stepsCompleted"`
		StepsTotal     int             `json:"totalSteps"`
		Data           json.RawMessage `json:"data"`
		Error          *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	var token string

	delay := 1 * time.Second
	const MAX_DELAY = 2 * time.Second
	const DELAY_FACTOR = 1.2

loop:
	for i := 0; ; i++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		respBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(respBytes, &respBody); err != nil {
			return nil, err
		}

		if !isSuccessful(resp) {
			errMsg := fmt.Sprintf("unsuccessful (%d) status code", resp.StatusCode)
			if respBody.Error != nil {
				errMsg += ": " + respBody.Error.Message
			}
			return nil, errors.New(errMsg)
		}

		// Only check for the token from the initial launch request
		if i == 0 {
			token = resp.Header.Get(TOKEN_HEADER)
		}

		if respBody.StepsCompleted >= respBody.StepsTotal {
			break
		}

		// Sleep but cancel if signaled by context
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(delay):
			// No-op
		}

		if delay < MAX_DELAY {
			delay = time.Duration(math.Round(float64(delay) * DELAY_FACTOR))
			if delay > MAX_DELAY {
				delay = MAX_DELAY
			}
		}

		u := fmt.Sprintf("%s/v2/api/queries/%s?lastStepSeen=%d", c.datasetUrl, respBody.Id, respBody.StepsCompleted)
		req, err = c.makeRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}
		if token != "" {
			req.Header.Set(TOKEN_HEADER, token)
		}
	}

	u := fmt.Sprintf("%s/v2/api/queries/%s?lastStepSeen=%d", c.datasetUrl, respBody.Id, respBody.StepsCompleted)
	req, err = c.makeRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set(TOKEN_HEADER, token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read and close the body so the transport can re-use a persistent tcp connection
	io.ReadAll(resp.Body)
	resp.Body.Close()

	return respBody.Data, nil
}
