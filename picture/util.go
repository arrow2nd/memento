package picture

import "strings"

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
