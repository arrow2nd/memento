package logparser

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ParseLog: 指定されたディレクトリ内の最新のVRChatログファイルを解析し、ワールド訪問履歴を取得する
func ParseLog(logDirPath string) ([]WorldVisit, error) {
	logFilePath := findRecentLogFilePath(logDirPath)
	if logFilePath == "" {
		return nil, errors.New("最新のログファイルが見つかりませんでした")
	}

	log.Println("最新のログファイルを取得: ", logFilePath)

	return parseWorldVisitsFromLog(logFilePath)
}

// findRecentLogFilePath: 指定されたディレクトリ内の最新のログファイルを探す
func findRecentLogFilePath(logDirPath string) string {
	files, err := os.ReadDir(logDirPath)
	if err != nil {
		return ""
	}

	var latestFile os.DirEntry
	var latestTime time.Time

	// 最新のログファイルを探す
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// ログファイルでないならスキップ
		if !isLogFile(file.Name()) {
			continue
		}

		// 作成日がより新しいなら更新
		if latestFile == nil || info.ModTime().After(latestTime) {
			latestFile = file
			latestTime = info.ModTime()
		}
	}

	// 見つからなかった
	if latestFile == nil {
		return ""
	}

	return filepath.Join(logDirPath, latestFile.Name())
}

