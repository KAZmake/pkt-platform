package onec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPClient calls the real 1С HTTP-service endpoints.
type HTTPClient struct {
	baseURL  string
	user     string
	password string
	http     *http.Client
}

func NewHTTPClient(baseURL, user, password string) *HTTPClient {
	return &HTTPClient{
		baseURL:  baseURL,
		user:     user,
		password: password,
		http:     &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *HTTPClient) GetLoans() ([]Loan, error) {
	return httpGet[[]Loan](c, "/loans")
}

func (c *HTTPClient) GetSchedule(loanOneCID string) ([]ScheduleItem, error) {
	return httpGet[[]ScheduleItem](c, "/loans/"+loanOneCID+"/schedule")
}

func (c *HTTPClient) GetDebts(loanOneCID string) ([]DebtItem, error) {
	return httpGet[[]DebtItem](c, "/loans/"+loanOneCID+"/debts")
}

func httpGet[T any](c *HTTPClient, path string) (T, error) {
	var zero T
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return zero, fmt.Errorf("build request: %w", err)
	}
	req.SetBasicAuth(c.user, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return zero, fmt.Errorf("http get %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("1C API %s returned %d", path, resp.StatusCode)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return zero, fmt.Errorf("decode %s: %w", path, err)
	}
	return result, nil
}
