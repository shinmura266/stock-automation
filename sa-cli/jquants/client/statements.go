package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"kabu-analysis/jquants"
)

// StatementsClient 財務情報API用クライアント
type StatementsClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewStatementsClient 新しい財務情報クライアントを作成
func NewStatementsClient(baseURL string, httpClient *http.Client) *StatementsClient {
	return &StatementsClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// GetStatements 財務情報を取得
func (c *StatementsClient) GetStatements(idToken, code, date string, debug bool) (*jquants.FinancialStatementsResponse, error) {
	return c.GetStatementsWithPagination(idToken, code, date, "", debug)
}

// GetStatementsWithPagination pagination_keyを指定して財務情報を取得
func (c *StatementsClient) GetStatementsWithPagination(idToken, code, date, paginationKey string, debug bool) (*jquants.FinancialStatementsResponse, error) {
	// パラメータの構築
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
	requestURL := fmt.Sprintf("%s/fins/statements", c.baseURL)
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}
	if debug {
		log.Printf("リクエストURL: %s", requestURL)
	}

	// リクエストの作成
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Authorization", "Bearer "+idToken)

	// リクエストの実行
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP リクエストエラー: %w", err)
	}
	defer resp.Body.Close()

	// ステータスコードのチェック
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	// レスポンスのデコード
	var financialStatementsResp jquants.FinancialStatementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&financialStatementsResp); err != nil {
		return nil, fmt.Errorf("レスポンスのデコードエラー: %w", err)
	}

	return &financialStatementsResp, nil
}

// GetAllStatements pagination_keyを自動処理してすべての財務情報を取得
func (c *StatementsClient) GetAllStatements(idToken, code, date string, debug bool) (*jquants.FinancialStatementsResponse, error) {
	var allStatements []jquants.FinancialStatement
	paginationKey := ""

	for {
		resp, err := c.GetStatementsWithPagination(idToken, code, date, paginationKey, debug)
		if err != nil {
			return nil, err
		}

		// 取得したデータを追加
		allStatements = append(allStatements, resp.Statements...)

		// pagination_keyが空なら終了
		if resp.PaginationKey == "" {
			break
		}

		// 次のページングキーを設定
		paginationKey = resp.PaginationKey

		if debug {
			log.Printf("次のページング処理: pagination_key=%s, 取得済み件数=%d", paginationKey, len(allStatements))
		}
	}

	if debug {
		log.Printf("全データ取得完了: 総件数=%d", len(allStatements))
	}

	return &jquants.FinancialStatementsResponse{
		Statements:    allStatements,
		PaginationKey: "", // 全データ取得完了なので空
	}, nil
}
