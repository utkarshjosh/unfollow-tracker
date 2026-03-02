package fetcher

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// InstagramScraper handles Instagram profile and follower scraping using web API
type InstagramScraper struct {
	client      *http.Client
	userAgent   string
	proxyPool   *ProxyPool
	sessionCookie string
	csrfToken   string
}

// NewInstagramScraper creates a new scraper with optional session cookie
func NewInstagramScraper(proxyPool *ProxyPool, sessionCookie string) *InstagramScraper {
	return &InstagramScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		proxyPool:     proxyPool,
		sessionCookie: sessionCookie,
		csrfToken:     extractCSRFToken(sessionCookie),
	}
}

// ProfileData represents scraped profile data
type ProfileData struct {
	Username      string
	UserID        string
	FollowerCount int
	IsPublic      bool
}

// FetchProfile fetches profile data including user ID and follower count
func (s *InstagramScraper) FetchProfile(ctx context.Context, username string) (*ProfileData, error) {
	// Clean username
	username = strings.TrimPrefix(username, "@")
	username = strings.TrimSpace(username)

	// First try: Use web profile page to extract data
	profileURL := fmt.Sprintf("https://www.instagram.com/%s/", username)

	data, err := s.makeRequest(ctx, profileURL, true)
	if err == nil {
		// Extract profile data from embedded JSON
		profile, parseErr := s.parseProfileData(data)
		if parseErr == nil {
			return profile, nil
		}
		// Web parsing failed, fall through to API fallback
	}
	// Fallback: Try API endpoint if web request or parsing fails
	return s.fetchProfileViaAPI(ctx, username)
}

// fetchProfileViaAPI uses Instagram's internal API to fetch profile data
func (s *InstagramScraper) fetchProfileViaAPI(ctx context.Context, username string) (*ProfileData, error) {
	if s.sessionCookie == "" {
		return nil, fmt.Errorf("session cookie required for API fetch")
	}

	apiURL := fmt.Sprintf("https://www.instagram.com/api/v1/users/web_profile_info/?username=%s", username)

	data, err := s.makeRequest(ctx, apiURL, false)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile via API: %w", err)
	}

	var result struct {
		Data struct {
			User struct {
				ID            string `json:"id"`
				Username      string `json:"username"`
				IsPrivate     bool   `json:"is_private"`
				FollowerCount int    `json:"edge_followed_by.count"`
				EdgeFollowedBy struct {
					Count int `json:"count"`
				} `json:"edge_followed_by"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse profile API response: %w", err)
	}

	user := result.Data.User
	return &ProfileData{
		Username:      user.Username,
		UserID:        user.ID,
		FollowerCount: user.EdgeFollowedBy.Count,
		IsPublic:      !user.IsPrivate,
	}, nil
}

// FetchFollowers fetches followers using Instagram's internal API
// Returns: follower hashes, next cursor (empty if done), error
func (s *InstagramScraper) FetchFollowers(ctx context.Context, userID string, cursor string) ([]string, string, error) {
	if s.sessionCookie == "" {
		return nil, "", fmt.Errorf("session cookie required to fetch followers")
	}

	// Instagram's friendships API endpoint
	apiURL := fmt.Sprintf("https://www.instagram.com/api/v1/friendships/%s/followers/", userID)

	// Build query parameters
	params := url.Values{}
	params.Set("count", "200") // Max per request
	params.Set("search_surface", "follow_list_page")

	if cursor != "" {
		params.Set("max_id", cursor)
	}

	fullURL := apiURL + "?" + params.Encode()

	data, err := s.makeRequest(ctx, fullURL, false)
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			return nil, "", fmt.Errorf("session expired or invalid")
		}
		if strings.Contains(err.Error(), "429") {
			return nil, "", fmt.Errorf("rate limited")
		}
		return nil, "", fmt.Errorf("failed to fetch followers: %w", err)
	}

	var result struct {
		Users []struct {
			Username string `json:"username"`
			ID       string `json:"pk"`
		} `json:"users"`
		NextMaxID   string `json:"next_max_id"`
		BigList     bool   `json:"big_list"`
		PageSize    int    `json:"page_size"`
		Status      string `json:"status"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, "", fmt.Errorf("failed to parse followers response: %w", err)
	}

	if result.Status != "ok" {
		return nil, "", fmt.Errorf("API returned non-ok status: %s", result.Status)
	}

	// Hash followers for privacy
	followers := make([]string, 0, len(result.Users))
	for _, user := range result.Users {
		hash := hashUsername(user.Username)
		followers = append(followers, hash)
	}

	// Return next cursor if there are more followers
	nextCursor := ""
	if result.BigList && result.NextMaxID != "" {
		nextCursor = result.NextMaxID
	}

	return followers, nextCursor, nil
}

// FetchAllFollowers fetches all followers for a user (handles pagination)
func (s *InstagramScraper) FetchAllFollowers(ctx context.Context, userID string, delay time.Duration) ([]string, error) {
	var allFollowers []string
	cursor := ""
	pageCount := 0

	for {
		select {
		case <-ctx.Done():
			return allFollowers, ctx.Err()
		default:
		}

		followers, nextCursor, err := s.FetchFollowers(ctx, userID, cursor)
		if err != nil {
			return allFollowers, fmt.Errorf("failed on page %d: %w", pageCount, err)
		}

		allFollowers = append(allFollowers, followers...)
		pageCount++

		// No more pages
		if nextCursor == "" {
			break
		}

		cursor = nextCursor

		// Respect rate limiting with delay
		if delay > 0 {
			time.Sleep(delay)
		}
	}

	return allFollowers, nil
}

// makeRequest makes an HTTP request with proper headers and session handling
func (s *InstagramScraper) makeRequest(ctx context.Context, requestURL string, isWebPage bool) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers to mimic browser
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	// Sec-Fetch-Site: "none" for web pages, "same-origin" for API requests
	if isWebPage {
		req.Header.Set("Sec-Fetch-Site", "none")
	} else {
		req.Header.Set("Sec-Fetch-Site", "same-origin")
	}
	req.Header.Set("Cache-Control", "max-age=0")

	// Add session cookie if available
	if s.sessionCookie != "" {
		req.Header.Set("Cookie", s.sessionCookie)

		// Add CSRF token for API requests
		if !isWebPage && s.csrfToken != "" {
			req.Header.Set("X-CSRFToken", s.csrfToken)
			req.Header.Set("X-IG-App-ID", "936619743392459") // Web app ID
			req.Header.Set("X-IG-WWW-Claim", "0")            // Required for API requests (instaloader pattern)
			req.Header.Set("X-Requested-With", "XMLHttpRequest")
			req.Header.Set("Referer", "https://www.instagram.com/")
		}
	}

	// Apply proxy if available
	if s.proxyPool != nil && len(s.proxyPool.proxies) > 0 {
		proxyURL := s.proxyPool.Next()
		if proxyURL != "" {
			parsedProxy, err := url.Parse(proxyURL)
			if err == nil {
				s.client.Transport = &http.Transport{
					Proxy: http.ProxyURL(parsedProxy),
				}
			}
		}
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		// 302 redirect to login page indicates session expired
		return nil, fmt.Errorf("302: session expired - redirect to login")
	}

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("401: unauthorized - session expired or invalid")
	}

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("429: rate limited")
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("404: user not found")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// parseProfileData extracts profile data from Instagram web page
func (s *InstagramScraper) parseProfileData(data []byte) (*ProfileData, error) {
	// Look for embedded JSON data in the page
	// Instagram embeds profile data in a <script> tag with type="application/ld+json"
	// or in window._sharedData

	content := string(data)

	// Try to find _sharedData
	startMarker := "window._sharedData = "
	startIdx := strings.Index(content, startMarker)
	if startIdx == -1 {
		return nil, fmt.Errorf("could not find shared data in page")
	}

	startIdx += len(startMarker)
	endIdx := strings.Index(content[startIdx:], ";</script>")
	if endIdx == -1 {
		return nil, fmt.Errorf("could not find end of shared data")
	}

	jsonData := content[startIdx : startIdx+endIdx]

	var sharedData struct {
		EntryData struct {
			ProfilePage []struct {
				GraphQL struct {
					User struct {
						ID            string `json:"id"`
						Username      string `json:"username"`
						IsPrivate     bool   `json:"is_private"`
						EdgeFollowedBy struct {
							Count int `json:"count"`
						} `json:"edge_followed_by"`
					} `json:"user"`
				} `json:"graphql"`
			} `json:"ProfilePage"`
		} `json:"entry_data"`
	}

	if err := json.Unmarshal([]byte(jsonData), &sharedData); err != nil {
		return nil, fmt.Errorf("failed to parse shared data: %w", err)
	}

	if len(sharedData.EntryData.ProfilePage) == 0 {
		return nil, fmt.Errorf("no profile page data found")
	}

	user := sharedData.EntryData.ProfilePage[0].GraphQL.User
	return &ProfileData{
		Username:      user.Username,
		UserID:        user.ID,
		FollowerCount: user.EdgeFollowedBy.Count,
		IsPublic:      !user.IsPrivate,
	}, nil
}

// hashUsername creates a SHA256 hash of username for privacy
func hashUsername(username string) string {
	h := sha256.New()
	h.Write([]byte(username))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// extractCSRFToken extracts CSRF token from session cookie
func extractCSRFToken(cookie string) string {
	// Look for csrftoken in cookie string
	parts := strings.Split(cookie, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "csrftoken=") {
			return strings.TrimPrefix(part, "csrftoken=")
		}
	}
	return ""
}

// ProxyPool manages rotating proxies
type ProxyPool struct {
	proxies []string
	index   int
}

// NewProxyPool creates a new proxy pool
func NewProxyPool(proxies []string) *ProxyPool {
	return &ProxyPool{
		proxies: proxies,
		index:   0,
	}
}

// Next returns the next proxy in rotation
func (p *ProxyPool) Next() string {
	if len(p.proxies) == 0 {
		return ""
	}

	proxy := p.proxies[p.index]
	p.index = (p.index + 1) % len(p.proxies)
	return proxy
}

// SetSessionCookie updates the session cookie for the scraper
func (s *InstagramScraper) SetSessionCookie(cookie string) {
	s.sessionCookie = cookie
	s.csrfToken = extractCSRFToken(cookie)
}

// IsAuthenticated returns true if the scraper has a valid session
func (s *InstagramScraper) IsAuthenticated() bool {
	return s.sessionCookie != ""
}

// RateLimiter implements simple rate limiting with exponential backoff
type RateLimiter struct {
	baseDelay    time.Duration
	maxDelay     time.Duration
	currentDelay time.Duration
	failures     int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(baseDelay, maxDelay time.Duration) *RateLimiter {
	return &RateLimiter{
		baseDelay:    baseDelay,
		maxDelay:     maxDelay,
		currentDelay: baseDelay,
		failures:     0,
	}
}

// Wait applies the current delay and increases it on failure
func (r *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(r.currentDelay):
		return nil
	}
}

// MarkSuccess resets the delay on success
func (r *RateLimiter) MarkSuccess() {
	r.currentDelay = r.baseDelay
	r.failures = 0
}

// MarkFailure increases delay on failure
func (r *RateLimiter) MarkFailure() {
	r.failures++
	r.currentDelay = r.currentDelay * 2
	if r.currentDelay > r.maxDelay {
		r.currentDelay = r.maxDelay
	}
}

// ParseFollowerCount parses follower count from various formats
func ParseFollowerCount(countStr string) (int, error) {
	// Remove commas and spaces
	countStr = strings.ReplaceAll(countStr, ",", "")
	countStr = strings.ReplaceAll(countStr, " ", "")

	// Handle K, M suffixes
	multiplier := 1
	if strings.HasSuffix(countStr, "K") {
		multiplier = 1000
		countStr = strings.TrimSuffix(countStr, "K")
	} else if strings.HasSuffix(countStr, "M") {
		multiplier = 1000000
		countStr = strings.TrimSuffix(countStr, "M")
	}

	// Parse float for cases like "1.5K"
	var count float64
	if _, err := fmt.Sscanf(countStr, "%f", &count); err != nil {
		return 0, err
	}

	return int(count * float64(multiplier)), nil
}

// ChunkFollowers splits followers into chunks for processing
func ChunkFollowers(followers []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(followers); i += chunkSize {
		end := i + chunkSize
		if end > len(followers) {
			end = len(followers)
		}
		chunks = append(chunks, followers[i:end])
	}
	return chunks
}

// calculateChunkCount determines how many chunks needed for follower count
func CalculateChunkCount(followerCount, chunkSize int) int {
	if followerCount <= 0 {
		return 1
	}
	chunks := followerCount / chunkSize
	if followerCount%chunkSize > 0 {
		chunks++
	}
	return chunks
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
