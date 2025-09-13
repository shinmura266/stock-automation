package main

import (
	"log"

	"kabu-analysis/cmd"

	"github.com/joho/godotenv"
)

func main() {
	// .envファイルから環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: .envファイルの読み込みに失敗しました: %v", err)
		log.Println("システム環境変数を使用します")
	}

	cmd.Execute()
}
