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

// StatementsClient 財務情報API用クライアント
type StatementsClient struct {
	baseURL    string
	interval   int
	httpClient *http.Client
}

// NewStatementsClient 新しい財務情報クライアントを作成
func NewStatementsClient(baseURL string, interval int, httpClient *http.Client) *StatementsClient {
	return &StatementsClient{
		baseURL:    baseURL,
		interval:   interval,
		httpClient: httpClient,
	}
}

// GetStatements 財務情報を取得
func (c *StatementsClient) GetStatements(idToken, code, date string) ([]schema.FinancialStatement, error) {
	// パラメータ組み立て
	params := url.Values{}
	if code != "" {
		params.Add("code", code)
	}
	if date != "" {
		params.Add("date", date)
	}

	var result []schema.FinancialStatement
	for {
		resp, err := c.requestStatements(idToken, params)
		if err != nil {
			return nil, err
		}

		result = append(result, resp.Statements...)

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

func (c *StatementsClient) requestStatements(idToken string, params url.Values) (*schema.FinancialStatementsResponse, error) {
	// URLの構築
	requestURL := fmt.Sprintf("%s/fins/statements", c.baseURL)
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+idToken)

	slog.Debug("Statementsリクエスト開始", "requestURL", requestURL)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ステータスコードエラー: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var result schema.FinancialStatementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	slog.Debug("Statementsリクエスト完了", "count", len(result.Statements), "pagination_key", result.PaginationKey)
	return &result, nil
}
