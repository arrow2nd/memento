package autostart

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// レジストリキーのパス
const registryPath = `Software\Microsoft\Windows\CurrentVersion\Run`

// SetAutoStart: 自動起動の状態を設定
func SetAutoStart(appName string, enable bool) error {
	// レジストリキーを開く
	key, err := registry.OpenKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	if err != nil {
		log.Println("レジストリキーのオープンに失敗:", err)
		return err
	}
	defer key.Close()

	// レジストリから削除
	if !enable {
		log.Println("自動起動を無効化")
		return key.DeleteValue(appName)
	}

	// 実行可能ファイルのパスを取得
	execPath, err := os.Executable()
	if err != nil {
		log.Println("実行ファイルパスの取得に失敗:", err)
		return err
	}

	execPath = filepath.Clean(execPath)

	// レジストリに登録
	if err := key.SetStringValue(appName, execPath); err != nil {
		log.Println("レジストリへの登録に失敗:", err)
		return err
	}

	log.Println("自動起動を有効化:", execPath)

	return nil
}

// IsAutoStartEnabled: 自動起動が有効かどうか
func IsAutoStartEnabled(appName string) bool {
	// レジストリキーを開く
	key, err := registry.OpenKey(registry.CURRENT_USER, registryPath, registry.QUERY_VALUE)
	if err != nil {
		log.Println("レジストリキーのオープンに失敗:", err)
		return false
	}
	defer key.Close()

	// レジストリ値を取得
	if _, _, err = key.GetStringValue(appName); err != nil {
		return false
	}

	return true
}
