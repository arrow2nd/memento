package logparser

import "regexp"

// isLogFile: ファイル名がログファイルの命名規則に従っているか
func isLogFile(fileName string) bool {
	pattern := `^output_log_\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}\.txt$`
	matched, err := regexp.MatchString(pattern, fileName)
	return err == nil && matched
}
