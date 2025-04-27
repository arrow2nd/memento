//go:build !windows

package autostart

import "log"

func SetAutoStart(appName string, enable bool) error {
	log.Println("SetAutoStart: 未対応のプラットフォーム")
	return nil
}

func IsAutoStartEnabled(appName string) bool {
	log.Println("IsAutoStartEnabled: 未対応のプラットフォーム")
	return false
}
