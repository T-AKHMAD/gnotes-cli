package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Note struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type NotesListResponse struct {
	Notes []Note `json:"notes"`
}

func (c *Client) NotesList(ctx context.Context, token string) (NotesListResponse, error) {
	url := c.baseURL + "/notes"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NotesListResponse{}, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return NotesListResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var ae apiError
		_ = json.NewDecoder(resp.Body).Decode(&ae)

		if ae.Error == "" {
			ae.Error = fmt.Sprintf("unexpected status: %s", resp.Status)
		}
		return NotesListResponse{}, fmt.Errorf(ae.Error)
	}

	var out NotesListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return NotesListResponse{}, err
	}

	return out, nil
}

type CreateNoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type CreateNoteResponse struct {
	ID int64 `json:"id"`
}

func (c *Client) CreateNote(ctx context.Context, token string, req CreateNoteRequest) (CreateNoteResponse, error) {
	url := c.baseURL + "/notes"

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return CreateNoteResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return CreateNoteResponse{}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return CreateNoteResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var ae apiError
		_ = json.NewDecoder(resp.Body).Decode(&ae)
		if ae.Error == "" {
			ae.Error = fmt.Sprintf("unexpected status: %s", resp.Status)
		}
		return CreateNoteResponse{}, fmt.Errorf(ae.Error)
	}

	var out CreateNoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return CreateNoteResponse{}, err
	}
	if out.ID <= 0 {
		return CreateNoteResponse{}, fmt.Errorf("unexpected create note response")
	}
	return out, nil
}

func (c *Client) GetNote(ctx context.Context, token string, id int64) (Note, error) {
	if id <= 0 {
		return Note{}, fmt.Errorf("invalid id")
	}

	url := fmt.Sprintf("%s/notes/%d", c.baseURL, id)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Note{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return Note{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var ae apiError
		_ = json.NewDecoder(resp.Body).Decode(&ae)
		if ae.Error == "" {
			ae.Error = fmt.Sprintf("unexpected status: %s", resp.Status)
		}
		return Note{}, fmt.Errorf(ae.Error)
	}

	var out Note
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Note{}, err
	}
	if out.ID <= 0 {
		return Note{}, fmt.Errorf("unexpected get note response")
	}
	return out, nil
}

func (c *Client) DeleteNote(ctx context.Context, token string, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid id")
	}

	url := fmt.Sprintf("%s/notes/%d", c.baseURL, id)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return nil
	}

	var ae apiError
	_ = json.NewDecoder(resp.Body).Decode(&ae)
	if ae.Error == "" {
		ae.Error = fmt.Sprintf("unexpected status: %s", resp.Status)
	}
	return fmt.Errorf(ae.Error)
}
