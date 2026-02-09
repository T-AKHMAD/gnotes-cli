package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type MeResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type apiError struct {
	Error string `json:"error"`
}

func (c *Client) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	url := c.baseURL + "/login"

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return LoginResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return LoginResponse{}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return LoginResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var ae apiError
		_ = json.NewDecoder(resp.Body).Decode(&ae)

		if ae.Error == "" {
			ae.Error = fmt.Sprintf("unexpected status: %s", resp.Status)
		}
		return LoginResponse{}, fmt.Errorf(ae.Error)
	}

	var out LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return LoginResponse{}, err
	}

	if out.Token == "" {
		return LoginResponse{}, fmt.Errorf("empty token in response")
	}

	return out, nil

}

func (c *Client) Me(ctx context.Context, token string) (MeResponse, error) {
	url := c.baseURL + "/me"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return MeResponse{}, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return MeResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var ae apiError
		_ = json.NewDecoder(resp.Body).Decode(&ae)

		if ae.Error == "" {
			ae.Error = fmt.Sprintf("unexpected status: %s", resp.Status)
		}
		return MeResponse{}, fmt.Errorf(ae.Error)
	}

	var out MeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return MeResponse{}, err
	}

	if out.ID == 0 || out.Email == "" {
		return MeResponse{}, fmt.Errorf("unexpected me response")
	}

	return out, nil
}
