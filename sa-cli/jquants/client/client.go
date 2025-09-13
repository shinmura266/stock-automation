package client

import (
	"net/http"
	"time"
)

// Client J-Quants APIの統合クライアント
type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *AuthClient
	listed     *ListedClient
	Daily      *DailyClient
	statements *StatementsClient
}

// NewClient 新しい統合クライアントを作成
func NewClient(baseURL string) *Client {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		auth:       NewAuthClient(baseURL, httpClient),
		listed:     NewListedClient(baseURL, httpClient),
		Daily:      NewDailyClient(baseURL, httpClient),
		statements: NewStatementsClient(baseURL, httpClient),
	}
}

// GetAuthClient 認証クライアントを取得
func (c *Client) GetAuthClient() *AuthClient {
	return c.auth
}

// GetListedClient 上場銘柄クライアントを取得
func (c *Client) GetListedClient() *ListedClient {
	return c.listed
}

// GetDailyClient 日次四本値クライアントを取得
func (c *Client) GetDailyClient() *DailyClient {
	return c.Daily
}

// GetStatementsClient 財務情報クライアントを取得
func (c *Client) GetStatementsClient() *StatementsClient {
	return c.statements
}

// Authenticate 認証を実行
func (c *Client) Authenticate(email, password string) error {
	// ステップ1: リフレッシュトークンを取得
	if err := c.auth.GetRefreshToken(email, password); err != nil {
		return err
	}

	// ステップ2: IDトークンを取得
	if err := c.auth.GetIDToken(); err != nil {
		return err
	}

	return nil
}

// GetIDToken IDトークンの値を取得
func (c *Client) GetIDToken() string {
	return c.auth.GetIDTokenValue()
}
