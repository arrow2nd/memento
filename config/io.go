package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// load: 設定を読み込む
func load(configPath string) (*Config, error) {
	buf, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗: %w", err)
	}

	// デフォルト設定
	config, err := getDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("デフォルト設定の取得に失敗: %w", err)
	}

	if err := json.Unmarshal(buf, &config); err != nil {
		return nil, fmt.Errorf("設定ファイルのデコードに失敗: %w", err)
	}

	config.configPath = configPath

	return config, nil
}

// Save: 設定を保存する
func (c *Config) Save() error {
	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(filepath.Dir(c.configPath), os.ModePerm); err != nil {
		return fmt.Errorf("設定ファイルのディレクトリ作成に失敗: %w", err)
	}

	// JSONにエンコード
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("設定ファイルのエンコードに失敗: %w", err)
	}

	if err := os.WriteFile(c.configPath, data, os.ModePerm); err != nil {
		return fmt.Errorf("設定ファイルの書き込みに失敗: %w", err)
	}

	return nil
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

// SetRootDirPath: 監視対象のルートディレクトリのパスを設定
func (c *Config) SetRootDirPath(path string) error {
	c.PictureDirPath = path
	return c.Save()
}

// SetVRCLogDirPath: VRChatのログディレクトリのパスを設定
func (c *Config) SetVRCLogDirPath(path string) error {
	c.VRCLogDirPath = path
	return c.Save()
}

// SetConvertToJpeg: JPEG変換の設定を変更
func (c *Config) SetConvertToJpeg(convert bool) error {
	c.ConvertToJpeg = convert
	return c.Save()
}

// SetJpegQuality: JPEG品質の設定を変更
func (c *Config) SetJpegQuality(quality int) error {
	// 品質は1-100の範囲内に制限
	if quality < 1 {
		quality = 1
	} else if quality > 100 {
		quality = 100
	}

	c.JpegQuality = quality
	return c.Save()
}
