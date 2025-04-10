package watcher

import (
	"log"
	"path/filepath"
)

// isCurrentMonthDir: 現在の年月のディレクトリかどうかを判定
func isCurrentMonthDir(path string) (string, bool) {
	// ディレクトリ名を取得
	dirName := filepath.Base(path)

	return dirName, dirName == getCurrentDate()
}

// isVRCPicture: VRChatの画像かどうかを判定
func (w *Watcher) isVRCPicture(path string) bool {
	// 親ディレクトリ名を取得
	dirName := filepath.Base(filepath.Dir(path))

	// 監視対象のディレクトリ配下でないなら無視
	if dirName != w.watchingSubDirName {
		log.Println("監視対象のディレクトリ配下でない: ", dirName)
		return false
	}

	ext := filepath.Ext(path)

	return ext == ".png" || ext == ".jpg"
}

