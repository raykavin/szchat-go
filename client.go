package szchat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/raykavin/gobox/httpclient"
	"github.com/raykavin/gobox/retry"
)

const (
	// apiVersionPath is the API version path this library targets. It is
	// appended automatically to baseURL by NewClient to keep internal
	// compatibility with this version of the SZChat API.
	apiVersionPath = "/api/v4"

	defaultMaxRetries   = 3
	defaultRetryWaitMin = 200 * time.Millisecond
	defaultRetryWaitMax = 2 * time.Second
)

// Client is the main SZChat SDK client. It is safe for concurrent use.
type Client struct {
	email       string
	password    string
	deviceToken string
	baseURL     string

	httpClient *http.Client

	tokenMu sync.RWMutex
	token   string

	maxRetries   int
	retryWaitMin time.Duration
	retryWaitMax time.Duration

	ContactAPI      *ContactAPI
	ContactGroupAPI *ContactGroupAPI
	ChannelAPI      *ChannelAPI
	AgentAPI        *AgentAPI
	AgentTalkAPI    *AgentTalkAPI
	AgentPauseAPI   *AgentPauseAPI
	TeamAPI         *TeamAPI
	AdminAPI        *AdminAPI
	MessageAPI      *MessageAPI
	ApplicationAPI  *ApplicationAPI
	TabulationAPI   *TabulationAPI
	PauseAPI        *PauseAPI
	AttendanceAPI   *AttendanceAPI
	TagAPI          *TagAPI
	HSMAPI          *HSMAPI
	GalleryAPI      *GalleryAPI
	PlaceholderAPI  *PlaceholderAPI
	CopilotAPI      *CopilotAPI
	WebRTCAPI       *WebRTCAPI
	TranslationAPI  *TranslationAPI
	ClickToCallAPI  *ClickToCallAPI
}

// Option is a function that configures a Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client used for all requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithDeviceToken sets the device token sent on login, used by SZChat for
// push notifications.
func WithDeviceToken(deviceToken string) Option {
	return func(c *Client) {
		c.deviceToken = deviceToken
	}
}

// WithRetry configures the retry policy applied to transient failures
// (429/502/503/504 and network errors).
func WithRetry(maxRetries int, waitMin, waitMax time.Duration) Option {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.retryWaitMin = waitMin
		c.retryWaitMax = waitMax
	}
}

// NewClient creates a new Client authenticating as the given agent/admin
// email and password against the SZChat API hosted at baseURL. baseURL is
// tenant-specific (e.g. "https://your-tenant.sz.chat") and must not include
// the API version path: NewClient appends "/api/v4" automatically to keep
// internal compatibility with the version of the SZChat API this library
// targets. There is no default baseURL. Options can override the HTTP
// client, device token, and retry policy.
func NewClient(baseURL, email, password string, opts ...Option) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("base url cannot be empty")
	}
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	c := &Client{
		email:        email,
		password:     password,
		baseURL:      strings.TrimRight(baseURL, "/") + apiVersionPath,
		httpClient:   http.DefaultClient,
		maxRetries:   defaultMaxRetries,
		retryWaitMin: defaultRetryWaitMin,
		retryWaitMax: defaultRetryWaitMax,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.ContactAPI = &ContactAPI{client: c}
	c.ContactGroupAPI = &ContactGroupAPI{client: c}
	c.ChannelAPI = &ChannelAPI{client: c}
	c.AgentAPI = &AgentAPI{client: c}
	c.AgentTalkAPI = &AgentTalkAPI{client: c}
	c.AgentPauseAPI = &AgentPauseAPI{client: c}
	c.TeamAPI = &TeamAPI{client: c}
	c.AdminAPI = &AdminAPI{client: c}
	c.MessageAPI = &MessageAPI{client: c}
	c.ApplicationAPI = &ApplicationAPI{client: c}
	c.TabulationAPI = &TabulationAPI{client: c}
	c.PauseAPI = &PauseAPI{client: c}
	c.AttendanceAPI = &AttendanceAPI{client: c}
	c.TagAPI = &TagAPI{client: c}
	c.HSMAPI = &HSMAPI{client: c}
	c.GalleryAPI = &GalleryAPI{client: c}
	c.PlaceholderAPI = &PlaceholderAPI{client: c}
	c.CopilotAPI = &CopilotAPI{client: c}
	c.WebRTCAPI = &WebRTCAPI{client: c}
	c.TranslationAPI = &TranslationAPI{client: c}
	c.ClickToCallAPI = &ClickToCallAPI{client: c}

	return c, nil
}

func (c *Client) setToken(token string) {
	c.tokenMu.Lock()
	c.token = token
	c.tokenMu.Unlock()
}

func (c *Client) getToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token
}

// noAuthRetryKey marks a context so that do() will not attempt an automatic
// token refresh on a 401 response. It is used internally by auth.go to avoid
// refreshToken recursively re-entering the authentication endpoints.
type noAuthRetryKey struct{}

func withNoAuthRetry(ctx context.Context) context.Context {
	return context.WithValue(ctx, noAuthRetryKey{}, true)
}

func hasNoAuthRetry(ctx context.Context) bool {
	skip, _ := ctx.Value(noAuthRetryKey{}).(bool)
	return skip
}

// statusError signals to retry.Do that the last attempt failed with a
// retryable (or refreshable) HTTP status and should be attempted again.
type statusError struct{ code int }

func (e *statusError) Error() string { return fmt.Sprintf("szchat: http status %d", e.code) }

func isRetryableStatus(code int) bool {
	switch code {
	case http.StatusTooManyRequests, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// rawResponse is a fully-buffered HTTP response.
type rawResponse struct {
	StatusCode  int
	Body        []byte
	ContentType string
}

// do executes req, transparently retrying on transient errors and on a
// single 401 after refreshing the bearer token.
func (c *Client) do(ctx context.Context, req *http.Request) (*rawResponse, error) {
	var bodyBytes []byte
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("szchat: reading request body: %w", err)
		}
		bodyBytes = b
	}

	var result *rawResponse
	refreshedOnce := false

	attempt := func() error {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			req.ContentLength = int64(len(bodyBytes))
		}

		if token := c.getToken(); token != "" {
			req.Header.Set(httpclient.HeaderAuthorization, "Bearer "+token)
		}
		if req.Header.Get(httpclient.HeaderAccept) == "" {
			req.Header.Set(httpclient.HeaderAccept, httpclient.MIMEApplicationJSON)
		}

		resp, err := c.httpClient.Do(req.Clone(ctx))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("szchat: reading response body: %w", err)
		}

		if resp.StatusCode == http.StatusUnauthorized && !refreshedOnce && !hasNoAuthRetry(ctx) {
			refreshedOnce = true
			if refreshErr := c.refreshToken(ctx); refreshErr == nil {
				return &statusError{code: resp.StatusCode}
			}
		}

		result = &rawResponse{
			StatusCode:  resp.StatusCode,
			Body:        respBody,
			ContentType: resp.Header.Get(httpclient.HeaderContentType),
		}
		if isRetryableStatus(resp.StatusCode) {
			return &statusError{code: resp.StatusCode}
		}
		return nil
	}

	shouldRetry := func(_ int, _ error) bool {
		return ctx.Err() == nil
	}

	if err := retry.Do(
		ctx,
		c.maxRetries,
		c.retryWaitMin,
		c.retryWaitMax,
		shouldRetry,
		attempt,
	); err != nil {
		var se *statusError
		if errors.As(err, &se) && result != nil {
			return result, nil
		}
		return nil, err
	}

	return result, nil
}

// request builds and executes a JSON HTTP request against path, encoding in
// as the request body (when non-nil) and decoding the response body into out
// (when non-nil). Non-2xx responses are returned as *APIError.
func (c *Client) request(
	ctx context.Context,
	method, path string,
	query url.Values,
	in, out any,
) error {
	return c.requestWithHeaders(ctx, method, path, query, nil, in, out)
}

// requestWithHeaders is like request but sets any extra headers on the
// request after the default Authorization/Accept/Content-Type headers,
// letting callers add or override headers (e.g. a channel-specific API key
// used instead of the agent bearer token) without duplicating the
// encode/decode/error-handling logic in request.
func (c *Client) requestWithHeaders(
	ctx context.Context,
	method, path string,
	query url.Values,
	headers map[string]string,
	in, out any,
) error {
	reqURL := c.baseURL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	var body io.Reader
	if in != nil {
		encoded, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("szchat: encoding request body: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return fmt.Errorf("szchat: building request: %w", err)
	}
	if body != nil {
		req.Header.Set(
			httpclient.HeaderContentType,
			httpclient.MIMEApplicationJSON,
		)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return parseAPIError(resp.StatusCode, resp.Body)
	}

	if out != nil && len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, out); err != nil {
			return fmt.Errorf("szchat: decoding response body: %w", err)
		}
	}

	return nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	return c.request(ctx, http.MethodGet, path, query, nil, out)
}

func (c *Client) post(ctx context.Context, path string, in, out any) error {
	return c.request(ctx, http.MethodPost, path, nil, in, out)
}

func (c *Client) put(ctx context.Context, path string, in, out any) error {
	return c.request(ctx, http.MethodPut, path, nil, in, out)
}

func (c *Client) delete(ctx context.Context, path string, out any) error {
	return c.request(ctx, http.MethodDelete, path, nil, nil, out)
}

// postWithKey issues a POST authenticated with a single static header
// (headerName: headerValue) instead of the client's bearer token, used by
// endpoints authenticated per-channel rather than per-agent (e.g. the
// generic channel and WhatsApp receptive APIs). The request still carries
// whatever bearer token the client currently holds, if any, since these
// endpoints only look at their own header and ignore Authorization.
func (c *Client) postWithKey(ctx context.Context, path, headerName, headerValue string, in, out any) error {
	return c.requestWithHeaders(ctx, http.MethodPost, path, nil, map[string]string{headerName: headerValue}, in, out)
}

// getRaw issues a GET and returns the raw response body and Content-Type
// header, for endpoints that do not return JSON (e.g. downloading stored
// media). Non-2xx responses are returned as *APIError.
func (c *Client) getRaw(ctx context.Context, path string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, "", fmt.Errorf("szchat: building request: %w", err)
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, "", parseAPIError(resp.StatusCode, resp.Body)
	}
	return resp.Body, resp.ContentType, nil
}

// postMultipart issues a multipart/form-data POST, used by endpoints that
// accept a file upload (e.g. the agent profile photo). fields holds
// additional plain form fields sent alongside the file.
func (c *Client) postMultipart(
	ctx context.Context,
	path, fileField, filename string,
	file io.Reader,
	fields map[string]string,
	out any,
) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fileField, filename)
	if err != nil {
		return fmt.Errorf("szchat: building multipart request: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("szchat: reading file content: %w", err)
	}
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return fmt.Errorf("szchat: building multipart request: %w", err)
		}
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("szchat: building multipart request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, &buf)
	if err != nil {
		return fmt.Errorf("szchat: building request: %w", err)
	}
	req.Header.Set(httpclient.HeaderContentType, writer.FormDataContentType())

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return parseAPIError(resp.StatusCode, resp.Body)
	}
	if out != nil && len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, out); err != nil {
			return fmt.Errorf("szchat: decoding response body: %w", err)
		}
	}
	return nil
}
