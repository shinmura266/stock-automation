package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"stock-automation/schema"
)

// ListedClient 上場銘柄関連のAPIクライアント
type ListedClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewListedClient 新しい上場銘柄クライアントを作成
func NewListedClient(baseURL string, httpClient *http.Client) *ListedClient {
	return &ListedClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// GetListedInfo 上場銘柄一覧を取得
func (c *ListedClient) GetListedInfo(idToken, date string) (*schema.ListedInfoResponse, error) {
	// パラメータ組み立て
	params := url.Values{}
	if date != "" {
		params.Add("date", date)
	}

	// URLの構築
	requestURL := fmt.Sprintf("%s/listed/info", c.baseURL)
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+idToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ステータスコードエラー: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var listedResp schema.ListedInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&listedResp); err != nil {
		return nil, err
	}

	return &listedResp, nil
}
