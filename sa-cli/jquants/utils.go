package jquants

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveToFile データをファイルに保存
func SaveToFile(data *ListedInfoResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return err
	}

	fmt.Printf("データを %s に保存しました\n", filename)
	return nil
}
