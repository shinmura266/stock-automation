package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"kabu-analysis/jquants"
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
func (c *ListedClient) GetListedInfo(idToken string) (*jquants.ListedInfoResponse, error) {
	return c.GetListedInfoWithPagination(idToken, "", "", "")
}

// GetListedInfoWithPagination pagination_keyを指定して上場銘柄一覧を取得
func (c *ListedClient) GetListedInfoWithPagination(idToken, code, date, paginationKey string) (*jquants.ListedInfoResponse, error) {
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

	var listedResp jquants.ListedInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&listedResp); err != nil {
		return nil, err
	}

	return &listedResp, nil
}

// GetListedInfoByCode 指定された銘柄コードの銘柄情報を取得
func (c *ListedClient) GetListedInfoByCode(idToken, code string) (*jquants.ListedInfoResponse, error) {
	return c.GetListedInfoWithPagination(idToken, code, "", "")
}

// GetListedInfoByDate 指定された日付の銘柄情報を取得
func (c *ListedClient) GetListedInfoByDate(idToken, date string) (*jquants.ListedInfoResponse, error) {
	return c.GetListedInfoWithPagination(idToken, "", date, "")
}

// GetAllListedInfo pagination_keyを自動処理してすべての上場銘柄一覧を取得
func (c *ListedClient) GetAllListedInfo(idToken, code, date string) (*jquants.ListedInfoResponse, error) {
	var allInfo []jquants.ListedInfo
	paginationKey := ""

	for {
		resp, err := c.GetListedInfoWithPagination(idToken, code, date, paginationKey)
		if err != nil {
			return nil, err
		}

		// 取得したデータを追加
		allInfo = append(allInfo, resp.Info...)

		// pagination_keyが空なら終了
		if resp.PaginationKey == "" {
			break
		}

		// 次のページングキーを設定
		paginationKey = resp.PaginationKey
	}

	return &jquants.ListedInfoResponse{
		Info:          allInfo,
		PaginationKey: "", // 全データ取得完了なので空
	}, nil
}
