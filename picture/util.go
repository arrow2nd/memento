package picture

import (
	"path/filepath"
	"strings"
)

// convertToSafeDirectoryName: ワールド名をディレクトリ名として使用可能な形に変換
func convertToSafeDirectoryName(worldName string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		".", "_",
	)

	return replacer.Replace(worldName)
}

// isMultiLayerPicture: マルチレイヤーの写真かどうかを判定
func isMultiLayerPicture(path string) bool {
	// 拡張子をチェック
	ext := filepath.Ext(path)
	if ext != ".png" && ext != ".jpg" {
		return false
	}

	// マルチレイヤーの写真ならtrue
	name := strings.TrimSuffix(filepath.Base(path), ext)
	return strings.HasSuffix(name, "_Environment") || strings.HasSuffix(name, "_Player")
}
