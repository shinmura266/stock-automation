package errors

import "fmt"

// 共通エラーメッセージ
const (
	ErrDatabaseConnection = "データベース接続エラー"
	ErrJQuantsAuth        = "J-Quants API認証エラー"
	ErrDataSave           = "データベース保存エラー"
	ErrDataRetrieve       = "データ取得エラー"
	ErrInvalidConfig      = "設定エラー"
)

// DatabaseError データベース関連のエラー
func DatabaseError(err error) error {
	return fmt.Errorf("%s: %v", ErrDatabaseConnection, err)
}

// JQuantsAuthError J-Quants認証関連のエラー
func JQuantsAuthError(err error) error {
	return fmt.Errorf("%s: %v", ErrJQuantsAuth, err)
}

// DataSaveError データ保存関連のエラー
func DataSaveError(err error) error {
	return fmt.Errorf("%s: %v", ErrDataSave, err)
}

// DataRetrieveError データ取得関連のエラー
func DataRetrieveError(err error) error {
	return fmt.Errorf("%s: %v", ErrDataRetrieve, err)
}

// ConfigError 設定関連のエラー
func ConfigError(err error) error {
	return fmt.Errorf("%s: %v", ErrInvalidConfig, err)
}
