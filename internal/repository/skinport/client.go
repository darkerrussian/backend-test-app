package skinport

import (
	"backend-test-app/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/andybalholm/brotli"
)

// Todo логичнее в будущем хранить в KV Consul, на случай если урл поменяется
const getItemsURL = "https://api.skinport.com/v1/items"

const (
	appIdURLKey    = "app_id"
	currencyURLKey = "currency"
	tradableURLKey = "tradable"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) GetItems(ctx context.Context, params entity.SkinportParams, tradable bool) ([]entity.SkinportRawItem, error) {
	query := url.Values{}
	query.Set(appIdURLKey, strconv.Itoa(params.AppId))
	query.Set(currencyURLKey, params.Currency)
	query.Set(tradableURLKey, strconv.FormatBool(tradable))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getItemsURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request error: %w", err)
	}
	// Важно: Accept-Encoding: br здесь обязателен, без него Skinport отвечает 406 Not Acceptable
	req.Header.Set("Accept-Encoding", "br")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "br" {
		reader = brotli.NewReader(resp.Body)
	}
	var items []entity.SkinportRawItem
	if err := json.NewDecoder(reader).Decode(&items); err != nil {
		return nil, fmt.Errorf("failed decode response: %w", err)
	}

	return items, nil
}
