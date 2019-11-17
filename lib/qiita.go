package lib

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Token      string
	Logger     *log.Logger
}

func New(rawBaseURL, token string, logger *log.Logger) (*Client, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = log.New(os.Stderr, "[LOG]", log.LstdFlags)
	}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		Token:      token,
		Logger:     logger,
	}, nil
}

type Item struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	LikesCount int    `json:"likes_count"`
}

func (c *Client) GetUserItems(ctx context.Context, userID string, page, perPage int) ([]Item, error) {
	reqURL := *c.BaseURL

	// set path
	reqURL.Path = path.Join(reqURL.Path, "users", userID, "items")

	// set query
	q := reqURL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("per_page", strconv.Itoa(perPage))
	reqURL.RawQuery = q.Encode()

	// instantiate request
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// set header
	req.Header.Set("User-Agent", "qiita-go-client")

	// set context
	req = req.WithContext(ctx)

	return nil, nil
}
