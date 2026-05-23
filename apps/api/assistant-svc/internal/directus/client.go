package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Program represents a loan program from Directus.
type Program struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	RateMin      float64 `json:"rate_min"`
	RateMax      float64 `json:"rate_max"`
	TermMin      int     `json:"term_min"`
	TermMax      int     `json:"term_max"`
	AmountMax    float64 `json:"amount_max"`
	Currency     string  `json:"currency"`
	Requirements string  `json:"requirements"`
}

// FAQ represents a FAQ item from Directus.
type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// Client fetches content from Directus REST API.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// GetPrograms fetches active loan programs.
func (c *Client) GetPrograms(ctx context.Context) ([]Program, error) {
	return directusGet[[]Program](ctx, c, "/items/programs?filter[status][_eq]=published&fields=name,description,rate_min,rate_max,term_min,term_max,amount_max,currency,requirements&limit=50")
}

// GetFAQ fetches published FAQ items.
func (c *Client) GetFAQ(ctx context.Context) ([]FAQ, error) {
	return directusGet[[]FAQ](ctx, c, "/items/faq?filter[status][_eq]=published&fields=question,answer&limit=100")
}

type directusResponse[T any] struct {
	Data T `json:"data"`
}

func directusGet[T any](ctx context.Context, c *Client, path string) (T, error) {
	var zero T
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return zero, fmt.Errorf("build request: %w", err)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return zero, fmt.Errorf("directus GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("directus %s returned %d", path, resp.StatusCode)
	}

	var envelope directusResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return zero, fmt.Errorf("decode %s: %w", path, err)
	}
	return envelope.Data, nil
}
