package picture

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// retryMoveFile: ファイル移動を複数回リトライする
func retryMoveFile(sourcePath, destPath string) {
	const (
		maxRetries = 5
		retryDelay = 3 * time.Second
	)

	for i := range maxRetries {
		// 少し待機
		time.Sleep(retryDelay)

		// 移動をリトライ
		err := os.Rename(sourcePath, destPath)
		if err == nil {
			log.Printf("リトライ成功 (%d回目): ファイルを移動: %s", i+1, destPath)
			return
		}

		log.Printf("リトライ失敗 (%d/%d): %v", i+1, maxRetries, err)
	}

	log.Printf("リトライ上限に達しました。ファイルの移動に失敗: %s, %s", sourcePath, destPath)
}

// createWorldNameDir: ワールド名のディレクトリを作成
func createWorldNameDir(targetDirPath, worldName string) error {
	// ワールド名をディレクトリ名に使用可能な形に変換
	safeWorldName := convertToSafeDirectoryName(worldName)
	worldDirPath := filepath.Join(targetDirPath, safeWorldName)

	// TargetDirPath以下にワールド名のディレクトリを作成
	if _, err := os.Stat(worldDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(worldDirPath, os.ModePerm); err != nil {
			return fmt.Errorf("ワールド名のディレクトリを作成できませんでした: %w", err)
		}
	}

	return nil
}
