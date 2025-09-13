package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"kabu-analysis/jquants"
)

// DailyClient 日次四本値関連のAPIクライアント
type DailyClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewDailyClient(baseURL string, httpClient *http.Client) *DailyClient {
	return &DailyClient{baseURL: baseURL, httpClient: httpClient}
}

// GetDailyQuotes 日次株価四本値を取得（code, date はいずれか/両方指定可）
func (c *DailyClient) GetDailyQuotes(idToken, code, date string) (*jquants.DailyQuotesResponse, error) {
	return c.GetDailyQuotesWithPagination(idToken, code, date, "")
}

// GetDailyQuotesWithPagination pagination_keyを指定して日次株価四本値を取得
func (c *DailyClient) GetDailyQuotesWithPagination(idToken, code, date, paginationKey string) (*jquants.DailyQuotesResponse, error) {
	// パラメータ組み立て
	params := url.Values{}
	if code != "" {
		params.Add("code", code)
	}
	if date != "" {
		params.Add("date", date)
	}
	if paginationKey != "" {
		params.Add("pagination_key", paginationKey)
	}

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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ステータスコードエラー: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var result jquants.DailyQuotesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllDailyQuotes pagination_keyを自動処理してすべての日次株価四本値を取得
func (c *DailyClient) GetAllDailyQuotes(idToken, code, date string) (*jquants.DailyQuotesResponse, error) {
	var allQuotes []jquants.DailyQuote
	paginationKey := ""

	for {
		resp, err := c.GetDailyQuotesWithPagination(idToken, code, date, paginationKey)
		if err != nil {
			return nil, err
		}

		// 取得したデータを追加
		allQuotes = append(allQuotes, resp.DailyQuotes...)

		// pagination_keyが空なら終了
		if resp.PaginationKey == "" {
			break
		}

		// 次のページングキーを設定
		paginationKey = resp.PaginationKey
	}

	return &jquants.DailyQuotesResponse{
		DailyQuotes:   allQuotes,
		PaginationKey: "", // 全データ取得完了なので空
	}, nil
}
