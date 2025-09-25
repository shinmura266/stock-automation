package cmd

import (
	"fmt"
	"log/slog"
	"stock-automation/jquants/service"

	"github.com/spf13/cobra"
)

var (
	listedInfoDate string
)

var ListedInfoCmd = &cobra.Command{
	Use:   "listed_info",
	Short: "上場銘柄一覧取得",
	Long:  "J-Quantsの上場銘柄一覧を取得して、DBへ保存する機能を提供します",
	RunE:  updateListedInfo,
}

func init() {
	// フラグを追加
	ListedInfoCmd.Flags().StringVarP(&listedInfoDate, "date", "d", "", "日付（YYYY-MM-DD形式、指定しない場合はAPIの最新日付）")
}

func updateListedInfo(cmd *cobra.Command, args []string) error {
	service, err := service.NewListedInfoService()
	if err != nil {
		return fmt.Errorf("上場銘柄情報サービス初期化エラー: %v", err)
	}
	defer service.Close()

	// 対象日付の全銘柄データを更新
	slog.Info("上場銘柄情報更新開始", "date", listedInfoDate)
	err = service.UpdateListedInfo(listedInfoDate)
	if err != nil {
		slog.Error("上場銘柄情報データ更新エラー", "error", err)
		return fmt.Errorf("上場銘柄情報データ更新エラー: %v", err)
	}
	slog.Info("上場銘柄情報データ更新完了")

	return nil
}
