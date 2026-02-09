package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Logout(ctx context.Context, token string) error {
	url := c.baseURL + "/logout"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	var ae apiError
	_ = json.NewDecoder(resp.Body).Decode(&ae)
	if ae.Error == "" {
		ae.Error = fmt.Sprintf("unexpected status: %s", resp.Status)
	}
	return fmt.Errorf(ae.Error)
}
