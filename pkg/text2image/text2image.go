package text2image

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const API_URL = "https://api.deepai.org/api/text2img"

type Client struct {
	apiKey string
	client *http.Client
}

type Response struct {
	ID        string `json:"id,omitempty"`
	OutputURL string `json:"output_url,omitempty"`
	Error     string `json:"err,omitempty"`
}

func NewClient() *Client {
	key := os.Getenv("DEEPAI_API_KEY")
	if key == "" {
		panic(errors.New("DEEPAI_API_KEY not defined"))
	}
	c := &http.Client{
		Timeout: 5 * time.Minute,
	}
	return &Client{apiKey: key, client: c}
}

func (c *Client) GetImage(text string) (string, error) {
	reader := bytes.NewBuffer([]byte(fmt.Sprintf("text=%s", text)))
	req, err := http.NewRequest(http.MethodPost, API_URL, reader)
	if err != nil {
		return "", fmt.Errorf("NewRequest: %w", err)
	}
	req.Header.Set("api-key", c.apiKey)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("io.ReadAll: %w", err)
	}

	r := &Response{}

	err = json.Unmarshal(body, r)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal: %w", err)
	}

	if r.Error != "" {
		return "", fmt.Errorf("error from api: %v", r)
	}
	if r.OutputURL == "" {
		return "", fmt.Errorf("empty output_url: %v", r)
	}

	return r.OutputURL, nil
}
