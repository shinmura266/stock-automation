package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"stock-automation/schema"
)

// AuthClient 認証関連のAPIクライアント
type AuthClient struct {
	baseURL          string
	httpClient       *http.Client
	accessTokenStore *AccessTokenStore
	tokenFilePath    string
}

// AccessTokenStore アクセストークンストア
type AccessTokenStore struct {
	MailAddress  string `json:"mailAddress"`
	RefreshToken *Token `json:"refreshToken"`
	IdToken      *Token `json:"idToken"`
}

// Token トークン情報
type Token struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// IsExpiredOrSoon トークンが期限切れか1時間以内に期限切れになるかどうかをチェック
func (t *Token) IsExpiredOrSoon() bool {
	if t == nil {
		return true
	}
	return time.Now().Add(time.Hour).After(t.ExpiresAt)
}

// invalidateTokens トークンを無効化する
func (c *AuthClient) invalidateTokens() {
	c.accessTokenStore.RefreshToken = nil
	c.accessTokenStore.IdToken = nil
	c.accessTokenStore.MailAddress = ""

	// ファイルに保存
	if err := c.saveTokensToFile(); err != nil {
		slog.Error("トークン無効化保存エラー", "error", err)
	}
}

// NewAuthClient 新しい認証クライアントを作成
func NewAuthClient(baseURL string, httpClient *http.Client) *AuthClient {
	// ホームディレクトリのパスを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// エラーの場合はデフォルトのパスを使用
		homeDir = "/home/ubuntu"
	}

	tokenFilePath := filepath.Join(homeDir, ".config", "stock-automation", "token")

	client := &AuthClient{
		baseURL:          baseURL,
		httpClient:       httpClient,
		tokenFilePath:    tokenFilePath,
		accessTokenStore: &AccessTokenStore{},
	}

	// トークンファイルが存在する場合は読み込む
	if err := client.loadTokensFromFile(); err != nil {
		// ファイルが存在しない場合は無視（初回実行時）
		if !os.IsNotExist(err) {
			slog.Error("トークンファイル読み込みエラー", "error", err)
		}
	}

	return client
}

// loadTokensFromFile トークンファイルからトークンを読み込む
func (c *AuthClient) loadTokensFromFile() error {
	data, err := os.ReadFile(c.tokenFilePath)
	if err != nil {
		return err
	}

	var store AccessTokenStore
	if err := json.Unmarshal(data, &store); err != nil {
		return fmt.Errorf("トークンファイルの解析エラー: %v", err)
	}

	c.accessTokenStore = &store
	return nil
}

// saveTokensToFile トークンをファイルに保存する
func (c *AuthClient) saveTokensToFile() error {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(c.tokenFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %v", err)
	}

	data, err := json.MarshalIndent(c.accessTokenStore, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %v", err)
	}

	if err := os.WriteFile(c.tokenFilePath, data, 0600); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %v", err)
	}

	return nil
}

// getRefreshToken リフレッシュトークンを取得（内部メソッド）
func (c *AuthClient) getRefreshToken() error {
	slog.Debug("getRefreshToken開始")

	// 環境変数から認証情報を取得
	mailAddress := os.Getenv("JQUANTS_EMAIL")
	password := os.Getenv("JQUANTS_PASSWORD")

	if mailAddress == "" || password == "" {
		return fmt.Errorf("JQUANTS_EMAILまたはJQUANTS_PASSWORD環境変数が設定されていません")
	}

	// メールアドレスが異なる場合は既存のトークンを無効化
	if c.accessTokenStore.MailAddress != "" && c.accessTokenStore.MailAddress != mailAddress {
		slog.Info("メールアドレスが変更されました。既存のトークンを無効化します", "old_email", c.accessTokenStore.MailAddress, "new_email", mailAddress)
		c.invalidateTokens()
	}

	// 既存のリフレッシュトークンが有効で、1時間以内に期限切れにならない場合はスキップ
	if c.accessTokenStore.RefreshToken != nil && !c.accessTokenStore.RefreshToken.IsExpiredOrSoon() {
		slog.Debug("リフレッシュトークンは有効です")
		return nil
	}

	url := fmt.Sprintf("%s/token/auth_user", c.baseURL)

	// リクエストボディを作成
	requestBody := schema.AuthUserRequest{
		Mailaddress: mailAddress,
		Password:    password,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		slog.Error("リクエストボディ作成エラー", "error", err)
		return fmt.Errorf("リクエストボディ作成エラー: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		slog.Error("HTTPリクエスト作成エラー", "error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("HTTPリクエスト実行エラー", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("ステータスコードエラー", "status_code", resp.StatusCode, "response", string(body))
		return fmt.Errorf("ステータスコードエラー: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var tokenResp schema.AuthUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		slog.Error("JSONデコードエラー", "error", err)
		return err
	}

	// リフレッシュトークンを保存（有効期限は1週間）
	c.accessTokenStore.MailAddress = mailAddress
	c.accessTokenStore.RefreshToken = &Token{
		Token:     tokenResp.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	// ファイルに保存
	if err := c.saveTokensToFile(); err != nil {
		slog.Error("トークン保存エラー", "error", err)
		return fmt.Errorf("トークン保存エラー: %v", err)
	}

	slog.Info("リフレッシュトークンを取得しました")
	return nil
}

// requestIdToken IDトークンを取得
func (c *AuthClient) requestIdToken() error {
	slog.Debug("requestIdToken開始", "idToken_exists", c.accessTokenStore.IdToken != nil)

	// 既存のIDトークンが有効で、1時間以内に期限切れにならない場合はスキップ
	if c.accessTokenStore.IdToken != nil && !c.accessTokenStore.IdToken.IsExpiredOrSoon() {
		slog.Debug("IDトークンは有効です")
		return nil
	}

	// リフレッシュトークンが存在しないか、有効期限が切れている場合は取得
	if c.accessTokenStore.RefreshToken == nil || c.accessTokenStore.RefreshToken.IsExpiredOrSoon() {
		if err := c.getRefreshToken(); err != nil {
			slog.Error("リフレッシュトークン取得エラー", "error", err)
			return fmt.Errorf("リフレッシュトークン取得エラー: %v", err)
		}
	}

	// クエリパラメータとしてリフレッシュトークンを設定
	url := fmt.Sprintf("%s/token/auth_refresh?refreshtoken=%s", c.baseURL, c.accessTokenStore.RefreshToken.Token)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		slog.Error("HTTPリクエスト作成エラー", "error", err)
		return err
	}

	// ヘッダーは不要（API仕様書に記載なし）

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("HTTPリクエスト実行エラー", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("ステータスコードエラー", "status_code", resp.StatusCode, "response", string(body))
		return fmt.Errorf("ステータスコードエラー: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var tokenResp schema.IdTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		slog.Error("JSONデコードエラー", "error", err)
		return err
	}

	// IDトークンを保存（有効期限は24時間）
	c.accessTokenStore.IdToken = &Token{
		Token:     tokenResp.IdToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// ファイルに保存
	if err := c.saveTokensToFile(); err != nil {
		slog.Error("トークン保存エラー", "error", err)
		return fmt.Errorf("トークン保存エラー: %v", err)
	}

	slog.Info("IDトークンを取得しました")
	return nil
}

// GetIdToken IDトークンの値を取得（有効期限切れの場合は自動取得）
func (c *AuthClient) GetIdToken() (string, error) {
	// IDトークンが存在しないか、有効期限が切れている場合は取得
	if c.accessTokenStore.IdToken == nil || c.accessTokenStore.IdToken.IsExpiredOrSoon() {
		slog.Debug("IDトークンが存在しないか期限切れのため、取得を開始します")
		if err := c.requestIdToken(); err != nil {
			slog.Error("IDトークン取得エラー", "error", err)
			return "", fmt.Errorf("IDトークン取得エラー: %v", err)
		}
	} else {
		slog.Debug("IDトークンは有効です", "expires_at", c.accessTokenStore.IdToken.ExpiresAt)
	}

	if c.accessTokenStore.IdToken == nil {
		slog.Error("IDトークンがnilです")
		return "", fmt.Errorf("IDトークンが取得できませんでした")
	}

	return c.accessTokenStore.IdToken.Token, nil
}
