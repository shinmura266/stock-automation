package cmd

import (
	"kabu-analysis/cmd/update"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "各データの更新を実行",
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(update.DailyCmd)
	updateCmd.AddCommand(update.DailyQuotesCmd)
	updateCmd.AddCommand(update.ListedInfoCmd)
	updateCmd.AddCommand(update.StatementsCmd)
	updateCmd.AddCommand(update.AssessmentCmd)
}
