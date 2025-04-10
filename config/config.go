package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	// RootDirPath: 監視対象のルートディレクトリのパス
	RootDirPath string
	// VRCLogDirPath: VRChatのログディレクトリのパス
	VRCLogDirPath string
}

// New: 新しいConfigを作成
func New() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
	}

	return &Config{
		RootDirPath:   getWatchDirPath(homeDir),
		VRCLogDirPath: getVRCLogDirPath(homeDir),
	}, nil
}

// getWatchDirPath: 監視対象のディレクトリのパスを取得
func getWatchDirPath(homeDir string) string {
	// TODO: 後で任意のディレクトリに変更できるようにする
	return filepath.Join(homeDir, "Pictures", "VRChat")
}

// getVRCLogDirPath: VRChatのログディレクトリのパスを取得
func getVRCLogDirPath(homeDir string) string {
	// TODO: 後で任意のディレクトリに変更できるようにする
	return filepath.Join(homeDir, "Documents", "VRChat")
}

