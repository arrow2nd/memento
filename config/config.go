package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// 設定ファイルの名前
const configFileName = "config.json"

type Config struct {
	// PictureDirPath: 写真の保存先のパス
	PictureDirPath string
	// VRCLogDirPath: VRChatのログディレクトリのパス
	VRCLogDirPath string

	configPath string
}

// New: 作成
func New(appName string) (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("Configディレクトリの取得に失敗: %w", err)
	}

	configPath := filepath.Join(configDir, appName, configFileName)
	log.Println("設定ファイル:", configPath)

	// 設定ファイルがあれば読み込む
	if _, err := os.Stat(configPath); err == nil {
		return load(configPath)
	}

	// デフォルト値を設定
	config := &Config{
		PictureDirPath: getDefaultWatchDirPath(homeDir),
		VRCLogDirPath:  getDefaultVRCLogDirPath(configDir),
		configPath:     configPath,
	}

	// 保存
	return config, config.Save()
}

// CheckDirectoriesExist: 設定されたディレクトリが存在するか確認する
func (c *Config) CheckDirectoriesExist() (bool, bool) {
	pictureExists, logExists := true, true

	// 写真ディレクトリの存在確認
	if _, err := os.Stat(c.PictureDirPath); err != nil {
		pictureExists = false
	}

	// ログディレクトリの存在確認
	if _, err := os.Stat(c.VRCLogDirPath); err != nil {
		logExists = false
	}

	return pictureExists, logExists
}

func getDefaultWatchDirPath(baseDir string) string {
	return filepath.Join(baseDir, "Pictures", "VRChat")
}

func getDefaultVRCLogDirPath(baseDir string) string {
	return filepath.Join(baseDir, "AppData", "LocalLow", "VRChat", "VRChat")
}
