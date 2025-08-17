package watcher

import (
	"log"
	"path/filepath"
	"strings"
)

// isCurrentMonthDir: 現在の年月のディレクトリかどうかを判定
func isCurrentMonthDir(path string) (string, bool) {
	// ディレクトリ名を取得
	dirName := filepath.Base(path)

	return dirName, dirName == getCurrentDate()
}

// isVRCPicture: VRChatの写真かどうかを判定
func (w *Watcher) isVRCPicture(path string) bool {
	// 親ディレクトリ名を取得
	dirName := filepath.Base(filepath.Dir(path))

	// 監視対象のディレクトリ配下でないなら無視
	if dirName != w.watchingSubDirName {
		log.Println("監視対象外なのでスキップ:", dirName)
		return false
	}

	filename := filepath.Base(path)
	ext := filepath.Ext(path)

	// マルチレイヤーの写真なら無視
	if name := strings.TrimSuffix(filename, ext); strings.HasSuffix(name, "_Environment") || strings.HasSuffix(name, "_Player") {
		log.Println("マルチレイヤーの写真なのでスキップ:", filename)
		return false
	}

	// 拡張子をチェック
	return ext == ".png" || ext == ".jpg"
}
