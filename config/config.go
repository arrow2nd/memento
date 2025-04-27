package config

import (
	"fmt"
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
	// ConvertToJpeg: JPEGに変換するかどうか
	ConvertToJpeg bool
	// JpegQuality: JPEGの品質設定 (1-100)
	JpegQuality int
	// ConfigDirPath: 設定ディレクトリのパス
	ConfigDirPath string

	// configFilePath: 設定ファイルのパス
	configFilePath string
}

// New: 作成
func New(appName string) (*Config, error) {
	config, err := getDefaultConfig(appName)
	if err != nil {
		return nil, fmt.Errorf("デフォルト設定の取得に失敗: %w", err)
	}

	// 設定ファイルがあれば読込む
	if _, err := os.Stat(config.configFilePath); err == nil {
		return load(config)
	}

	return config, config.Save()
}

// getConfigDirPath: 設定ディレクトリのパスを取得
func getConfigDirPath(appName string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, appName), nil
}

// getDefaultConfig: デフォルトの設定を取得
func getDefaultConfig(appName string) (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
	}

	configDir, err := getConfigDirPath(appName)
	if err != nil {
		return nil, fmt.Errorf("設定ディレクトリの取得に失敗: %w", err)
	}

	return &Config{
		ConfigDirPath:  configDir,
		PictureDirPath: getDefaultWatchDirPath(homeDir),
		VRCLogDirPath:  getDefaultVRCLogDirPath(homeDir),
		ConvertToJpeg:  true,
		JpegQuality:    90,
		configFilePath: filepath.Join(configDir, configFileName),
	}, nil
}

// getDefaultWatchDirPath: デフォルトの監視ディレクトリのパスを取得
func getDefaultWatchDirPath(baseDir string) string {
	return filepath.Join(baseDir, "Pictures", "VRChat") // OneDrive入ってたらなんか違うかも
}

// getDefaultVRCLogDirPath: デフォルトのVRChatログディレクトリのパスを取得
func getDefaultVRCLogDirPath(baseDir string) string {
	return filepath.Join(baseDir, "AppData", "LocalLow", "VRChat", "VRChat")
}
