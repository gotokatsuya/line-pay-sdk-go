package linepay

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

// API endpoint base constants
const (
	APIEndpointReal    = "https://api-pay.line.me"
	APIEndpointSandbox = "https://sandbox-api-pay.line.me"
)

// Client type
type Client struct {
	channelID     string
	channelSecret string
	endpoint      *url.URL
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
	if c.endpoint == nil {
		u, err := url.Parse(APIEndpointReal)
		if err != nil {
			return nil, err
		}
		c.endpoint = u
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

// WithEndpoint function
func WithEndpoint(endpoint string) ClientOption {
	return func(client *Client) error {
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		client.endpoint = u
		return nil
	}
}

// WithSandbox function
func WithSandbox() ClientOption {
	return WithEndpoint(APIEndpointSandbox)
}

// mergeQuery method
func (c *Client) mergeQuery(path string, q interface{}) (string, error) {
	v := reflect.ValueOf(q)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return path, nil
	}

	u, err := url.Parse(path)
	if err != nil {
		return path, err
	}

	qs, err := query.Values(q)
	if err != nil {
		return path, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewRequest method
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	message := c.channelSecret + path

	switch method {
	case http.MethodGet, http.MethodDelete:
		if body != nil {
			merged, err := c.mergeQuery(path, body)
			if err != nil {
				return nil, err
			}
			path = merged
		}
	}
	u, err := c.endpoint.Parse(path)
	if err != nil {
		return nil, err
	}

	var reqBody io.ReadWriter
	switch method {
	case http.MethodGet, http.MethodDelete:
		if body != nil {
			message += u.RawQuery
		}
	case http.MethodPost, http.MethodPut:
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewBuffer(b)
			message += string(b)
		}
	}

	nounce := uuid.New().String()
	message += nounce

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-LINE-ChannelId", c.channelID)
	req.Header.Set("X-LINE-Authorization-Nonce", nounce)

	hash := hmac.New(sha256.New, []byte(c.channelSecret))
	hash.Write([]byte(message))
	req.Header.Set("X-LINE-Authorization", base64.StdEncoding.EncodeToString(hash.Sum(nil)))
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
		return nil, err
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
