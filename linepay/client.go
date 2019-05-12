package linepay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// API endpoint base constants
const (
	APIEndpointBaseReal    = "https://api-pay.line.me/"
	APIEndpointBaseSandbox = "https://sandbox-api-pay.line.me/"
)

// Client type
type Client struct {
	channelID     string
	channelSecret string
	endpointBase  *url.URL
	httpClient    *http.Client
}

// ClientOption type
type ClientOption func(*Client) error

// New returns a new pay client instance.
func New(channelID, channelSecret string, options ...ClientOption) (*Client, error) {
	if channelID == "" {
		return nil, errors.New("missing channel id")
	}
	if channelSecret == "" {
		return nil, errors.New("missing channel secret")
	}
	c := &Client{
		channelID:     channelID,
		channelSecret: channelSecret,
		httpClient:    http.DefaultClient,
	}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}
	if c.endpointBase == nil {
		u, err := url.Parse(APIEndpointBaseReal)
		if err != nil {
			return nil, err
		}
		c.endpointBase = u
	}
	return c, nil
}

// WithHTTPClient function
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		client.httpClient = c
		return nil
	}
}

// WithEndpointBase function
func WithEndpointBase(endpointBase string) ClientOption {
	return func(client *Client) error {
		u, err := url.Parse(endpointBase)
		if err != nil {
			return err
		}
		client.endpointBase = u
		return nil
	}
}

// WithSandbox function
func WithSandbox() ClientOption {
	return WithEndpointBase(APIEndpointBaseSandbox)
}

// mergeQuery method
func mergeQuery(endpoint string, q interface{}) (string, error) {
	v := reflect.ValueOf(q)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return endpoint, nil
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return endpoint, err
	}

	qs, err := query.Values(q)
	if err != nil {
		return endpoint, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewRequest method
func (c *Client) NewRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	u, err := c.endpointBase.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-LINE-ChannelId", c.channelID)
	req.Header.Set("X-LINE-ChannelSecret", c.channelSecret)
	return req, nil
}

// Do method
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	defer resp.Body.Close()

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
				return resp, err
			}
		}
	}
	return resp, err
}
