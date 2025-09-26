package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"stock-automation/schema"
)

// DailyQuotesClient 日次株価四本値関連のAPIクライアント
type DailyQuotesClient struct {
	baseURL    string
	interval   int
	httpClient *http.Client
}

// NewDailyQuotesClient 新しい日次株価四本値クライアントを作成
func NewDailyQuotesClient(baseURL string, interval int, httpClient *http.Client) *DailyQuotesClient {
	return &DailyQuotesClient{
		baseURL:    baseURL,
		interval:   interval,
		httpClient: httpClient,
	}
}

// GetDailyQuotes 日次株価四本値を取得
func (c *DailyQuotesClient) GetDailyQuotes(idToken, code, date string) ([]schema.DailyQuote, error) {
	// パラメータ組み立て
	params := url.Values{}
	if code != "" {
		params.Add("code", code)
	}
	if date != "" {
		params.Add("date", date)
	}

	var result []schema.DailyQuote
	for {
		resp, err := c.requestDailyQuotes(idToken, params)
		if err != nil {
			return nil, err
		}

		result = append(result, resp.DailyQuotes...)

		if resp.PaginationKey == "" {
			break
		}

		params.Add("pagination_key", resp.PaginationKey)

		// PaginationKeyによる繰り返し時にintervalのインターバル
		if c.interval > 0 {
			time.Sleep(time.Duration(c.interval) * time.Second)
		}
	}

	return result, nil
}

func (c *DailyQuotesClient) requestDailyQuotes(idToken string, params url.Values) (*schema.DailyQuotesResponse, error) {
	// URLの構築
	requestURL := fmt.Sprintf("%s/prices/daily_quotes", c.baseURL)
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+idToken)

	slog.Debug("DailyQuotesリクエスト開始", "requestURL", requestURL)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ステータスコードエラー: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var result schema.DailyQuotesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	slog.Debug("DailyQuotesリクエスト完了", "count", len(result.DailyQuotes), "pagination_key", result.PaginationKey)
	return &result, nil
}
