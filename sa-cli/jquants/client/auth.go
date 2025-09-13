package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AuthClient 認証関連のAPIクライアント
type AuthClient struct {
	baseURL      string
	httpClient   *http.Client
	refreshToken string
	idToken      string
}

// NewAuthClient 新しい認証クライアントを作成
func NewAuthClient(baseURL string, httpClient *http.Client) *AuthClient {
	return &AuthClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// GetRefreshToken リフレッシュトークンを取得
func (c *AuthClient) GetRefreshToken(email, password string) error {
	url := fmt.Sprintf("%s/token/auth_user", c.baseURL)

	// リクエストボディを作成
	requestBody := map[string]string{
		"mailaddress": email,
		"password":    password,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("リクエストボディ作成エラー: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ステータスコードエラー: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	c.refreshToken = tokenResp.RefreshToken
	fmt.Println("リフレッシュトークンを取得しました")
	return nil
}

// GetIDToken IDトークンを取得
func (c *AuthClient) GetIDToken() error {
	if c.refreshToken == "" {
		return fmt.Errorf("リフレッシュトークンが設定されていません")
	}

	// クエリパラメータとしてリフレッシュトークンを設定
	url := fmt.Sprintf("%s/token/auth_refresh?refreshtoken=%s", c.baseURL, c.refreshToken)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	// ヘッダーは不要（API仕様書に記載なし）

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ステータスコードエラー: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		IDToken string `json:"idToken"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	c.idToken = tokenResp.IDToken
	fmt.Println("IDトークンを取得しました")
	return nil
}

// GetIDTokenValue IDトークンの値を取得
func (c *AuthClient) GetIDTokenValue() string {
	return c.idToken
}

// GetRefreshTokenValue リフレッシュトークンの値を取得
func (c *AuthClient) GetRefreshTokenValue() string {
	return c.refreshToken
}
