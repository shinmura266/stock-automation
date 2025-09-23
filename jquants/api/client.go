package api

import "net/http"

const baseURL = "https://api.jquants.com/v1"
const interval = 1000

// Client J-Quants APIクライアント
type Client struct {
	AuthClient        *AuthClient
	ListedClient      *ListedClient
	DailyQuotesClient *DailyQuotesClient
	StatementsClient  *StatementsClient
}

// NewClient 新しいクライアントを作成
func NewClient() *Client {
	httpClient := &http.Client{}
	return &Client{
		AuthClient:        NewAuthClient(baseURL, httpClient),
		ListedClient:      NewListedClient(baseURL, httpClient),
		DailyQuotesClient: NewDailyQuotesClient(baseURL, interval, httpClient),
		StatementsClient:  NewStatementsClient(baseURL, interval, httpClient),
	}
}
