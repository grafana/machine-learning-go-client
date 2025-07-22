package mlapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

type responseWrapper[T any] struct {
	Status   string   `json:"status"`
	Data     T        `json:"data"`
	Warnings []string `json:"warnings"`
	Error    string   `json:"error"`
}

// Client is a Grafana API client.
type Client struct {
	config  Config
	baseURL url.URL
	client  *http.Client
}

// Config contains client configuration.
type Config struct {
	// BearerToken is an optional API key.
	BearerToken string
	// BasicAuth is optional basic auth credentials.
	BasicAuth *url.Userinfo
	// Client provides an optional HTTP client, otherwise a default will be used.
	Client *http.Client
	// NumRetries contains the number of attempted retries
	NumRetries int
}

// New creates a new Grafana client.
func New(baseURL string, cfg Config) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if cfg.BasicAuth != nil {
		u.User = cfg.BasicAuth
	}

	cli := cfg.Client
	if cli == nil {
		cli = cleanhttp.DefaultClient()
	}

	return &Client{
		config:  cfg,
		baseURL: *u,
		client:  cli,
	}, nil
}

func (c *Client) request(ctx context.Context, method, requestPath string, query url.Values, body io.Reader, responseStruct any) error {
	var (
		req          *http.Request
		resp         *http.Response
		err          error
		bodyContents []byte
	)

	// read the request body and save it so we can use it in retries.
	var reqBody []byte
	if body != nil {
		reqBody, err = io.ReadAll(body)
		if err != nil {
			return fmt.Errorf("failed to read request body: %w", err)
		}
	}

	// retry logic
	for n := 0; n <= c.config.NumRetries; n++ {
		var bodyReader io.Reader
		if reqBody != nil {
			bodyReader = bytes.NewReader(reqBody)
		}
		req, err = c.newRequest(ctx, method, requestPath, query, bodyReader)
		if err != nil {
			return err
		}

		// Wait a bit if that's not the first request
		if n != 0 {
			time.Sleep(time.Second * 5)
		}

		resp, err = c.client.Do(req)

		// If err is not nil, retry again
		// That's either caused by client policy, or failure to speak HTTP (such as network connectivity problem). A
		// non-2xx status code doesn't cause an error.
		if err != nil {
			continue
		}

		//nolint:errcheck // We can't do anything about not being able to close the body.
		defer resp.Body.Close()

		// read the body (even on non-successful HTTP status codes), as that's what the unit tests expect
		bodyContents, err = io.ReadAll(resp.Body)

		// if there was an error reading the body, try again
		if err != nil {
			continue
		}

		// Exit the loop if we have something final to return. This is anything < 500, if it's not a 429.
		if resp.StatusCode < http.StatusInternalServerError && resp.StatusCode != http.StatusTooManyRequests {
			break
		}
	}
	if err != nil {
		return err
	}

	// check status code.
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status: %d, body: %v", resp.StatusCode, string(bodyContents))
	}

	if responseStruct == nil {
		return nil
	}

	err = json.Unmarshal(bodyContents, responseStruct)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(ctx context.Context, method, requestPath string, query url.Values, body io.Reader) (*http.Request, error) {
	url := c.baseURL
	url.Path = path.Join(url.Path, requestPath)
	url.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return req, err
	}

	if c.config.BearerToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.BearerToken))
	}

	req.Header.Add("Content-Type", "application/json")
	return req, err
}
