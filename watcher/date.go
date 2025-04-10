package watcher

import "time"

// getCurrentDate: 現在の年月を取得
func getCurrentDate() string {
	return time.Now().Format("2006-01")
}

