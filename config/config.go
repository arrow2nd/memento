package config

import (
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	RootDirPath string
}

// New: 新しいConfigを作成
func New() *Config {
	return &Config{
		RootDirPath: getWatchDirPath(),
	}
}

// getWatchDirPath: 監視対象のディレクトリのパスを取得
func getWatchDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("ホームディレクトリの取得に失敗: ", err)
	}

	// TODO: 後で任意のディレクトリに変更できるようにする
	return filepath.Join(homeDir, "Pictures", "VRChat")
}

