package cmd

import (
	"kabu-analysis/cmd/analyze"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "各種データの分析を実行",
	Long:  "保存された各種データを分析し、レポートやサマリーを作成します",
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.AddCommand(analyze.AnalyzeStatementsCmd)
	analyzeCmd.AddCommand(analyze.QueryCmd)
	analyzeCmd.AddCommand(analyze.AssessmentCmd)
}
