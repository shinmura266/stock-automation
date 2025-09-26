package cmd

import (
	"fmt"
	"log/slog"
	"stock-automation/jquants/service"
	"time"

	"github.com/spf13/cobra"
)

var (
	dailyDate     string
	dailyCount    int
	dailyInterval int
)

var DailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "日次データ一括更新",
	Long:  "上場銘柄一覧→日次株価四本値→財務情報の順で一括更新します",
	RunE:  updateDaily,
}

func init() {
	// フラグを追加
	DailyCmd.Flags().StringVarP(&dailyDate, "date", "d", "", "日付（YYYY-MM-DD形式、指定しない場合はAPIの最新日付）")
	DailyCmd.Flags().IntVarP(&dailyCount, "count", "c", 1, "取得する日数（指定した日付からさかのぼる日数、デフォルト: 1）")
	DailyCmd.Flags().IntVar(&dailyInterval, "interval", 5, "インターバル（秒、デフォルト: 5）")
}

func updateDaily(cmd *cobra.Command, args []string) error {
	// グローバルフラグからverboseの値を取得
	verbose, _ := cmd.Root().PersistentFlags().GetBool("verbose")

	slog.Info("日次データ一括更新開始", "date", dailyDate, "count", dailyCount)

	// 1. 上場銘柄一覧の更新
	slog.Info("1. 上場銘柄一覧更新開始")
	listedInfoService, err := service.NewListedInfoService(verbose)
	if err != nil {
		return fmt.Errorf("上場銘柄情報サービス初期化エラー: %v", err)
	}
	defer listedInfoService.Close()

	err = listedInfoService.UpdateListedInfo(dailyDate)
	if err != nil {
		slog.Error("上場銘柄情報データ更新エラー", "error", err)
		return fmt.Errorf("上場銘柄情報データ更新エラー: %v", err)
	}
	slog.Info("上場銘柄一覧更新完了")

	// インターバル待機
	slog.Debug("インターバル待機中", "interval", dailyInterval)
	time.Sleep(time.Duration(dailyInterval) * time.Second)

	// 2. 日次株価四本値の更新
	slog.Info("2. 日次株価四本値更新開始")
	dailyQuotesService, err := service.NewDailyQuotesService(dailyInterval, verbose)
	if err != nil {
		return fmt.Errorf("株価サービス初期化エラー: %v", err)
	}

	err = dailyQuotesService.UpdateDailyQuotesWithCount("", dailyDate, dailyCount)
	if err != nil {
		slog.Error("日次株価四本値データ更新エラー", "error", err)
		return fmt.Errorf("日次株価四本値データ更新エラー: %v", err)
	}
	slog.Info("日次株価四本値更新完了")

	// インターバル待機
	slog.Debug("インターバル待機中", "interval", dailyInterval)
	time.Sleep(time.Duration(dailyInterval) * time.Second)

	// 3. 財務情報の更新
	slog.Info("3. 財務情報更新開始")
	statementsService, err := service.NewStatementsService(dailyInterval, verbose)
	if err != nil {
		return fmt.Errorf("財務情報サービス初期化エラー: %v", err)
	}
	defer statementsService.Close()

	err = statementsService.UpdateStatementsWithCount("", dailyDate, dailyCount)
	if err != nil {
		slog.Error("財務情報データ更新エラー", "error", err)
		return fmt.Errorf("財務情報データ更新エラー: %v", err)
	}
	slog.Info("財務情報更新完了")

	slog.Info("日次データ一括更新完了")
	return nil
}
