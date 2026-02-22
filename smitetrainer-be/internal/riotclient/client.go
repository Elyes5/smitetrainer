package riotclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"smitetrainer-be/internal/cache"
)

type Client struct {
	apiKey      string
	httpClient  *http.Client
	maxAttempts int
	cache       cache.Store
	cacheTTL    time.Duration
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("riot api status %d: %s", e.StatusCode, e.Message)
}

func New(apiKey string, timeout time.Duration, maxAttempts int, responseCache cache.Store, cacheTTL time.Duration) *Client {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxAttempts: maxAttempts,
		cache:       responseCache,
		cacheTTL:    cacheTTL,
	}
}

func (c *Client) GetMatch(ctx context.Context, matchID string) (MatchResponse, error) {
	route, err := RegionalRouteFromMatchID(matchID)
	if err != nil {
		return MatchResponse{}, err
	}
	endpoint := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s", route, url.PathEscape(matchID))

	var out MatchResponse
	if err := c.fetchJSON(ctx, "match:"+matchID, endpoint, &out); err != nil {
		return MatchResponse{}, err
	}
	return out, nil
}

func (c *Client) GetTimeline(ctx context.Context, matchID string) (TimelineResponse, error) {
	route, err := RegionalRouteFromMatchID(matchID)
	if err != nil {
		return TimelineResponse{}, err
	}
	endpoint := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s/timeline", route, url.PathEscape(matchID))

	var out TimelineResponse
	if err := c.fetchJSON(ctx, "timeline:"+matchID, endpoint, &out); err != nil {
		return TimelineResponse{}, err
	}
	return out, nil
}

func (c *Client) fetchJSON(ctx context.Context, cacheKey, endpoint string, out any) error {
	if c.cache != nil {
		if cached, found, err := c.cache.Get(ctx, cacheKey); err == nil && found {
			if err := json.Unmarshal(cached, out); err == nil {
				return nil
			}
		}
	}

	body, err := c.doWithRetry(ctx, endpoint)
	if err != nil {
		return err
	}

	if c.cache != nil {
		_ = c.cache.Set(ctx, cacheKey, body, c.cacheTTL)
	}
	return json.Unmarshal(body, out)
}

func (c *Client) doWithRetry(ctx context.Context, endpoint string) ([]byte, error) {
	var lastErr error

	for attempt := 1; attempt <= c.maxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Riot-Token", c.apiKey)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request error: %w", err)
			if attempt == c.maxAttempts {
				break
			}
			if err := sleepWithContext(ctx, backoffDuration(attempt, "")); err != nil {
				return nil, err
			}
			continue
		}

		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
		_ = resp.Body.Close()
		if readErr != nil {
			lastErr = fmt.Errorf("failed to read riot response: %w", readErr)
			if attempt == c.maxAttempts {
				break
			}
			if err := sleepWithContext(ctx, backoffDuration(attempt, "")); err != nil {
				return nil, err
			}
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return body, nil
		}

		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Message:    trimBody(body),
		}
		lastErr = apiErr

		if !isRetryableStatus(resp.StatusCode) || attempt == c.maxAttempts {
			return nil, apiErr
		}

		if err := sleepWithContext(ctx, backoffDuration(attempt, resp.Header.Get("Retry-After"))); err != nil {
			return nil, err
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("riot request failed for %s", endpoint)
	}
	return nil, lastErr
}

func isRetryableStatus(code int) bool {
	return code == http.StatusTooManyRequests || code >= http.StatusInternalServerError
}

func backoffDuration(attempt int, retryAfterHeader string) time.Duration {
	const (
		baseDelay = 300 * time.Millisecond
		maxDelay  = 6 * time.Second
	)

	delay := baseDelay << (attempt - 1)
	if delay > maxDelay {
		delay = maxDelay
	}

	if retryAfterHeader != "" {
		if seconds, err := strconv.Atoi(strings.TrimSpace(retryAfterHeader)); err == nil && seconds > 0 {
			retryDelay := time.Duration(seconds) * time.Second
			if retryDelay > delay {
				delay = retryDelay
			}
		}
	}
	return delay
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func trimBody(body []byte) string {
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return http.StatusText(http.StatusBadGateway)
	}
	const maxLen = 240
	if len(trimmed) > maxLen {
		return trimmed[:maxLen] + "..."
	}
	return trimmed
}
