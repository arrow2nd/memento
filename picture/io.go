package picture

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	MAX_RETRIES = 5
	RETRY_DELAY = 3 * time.Second
)

// removeFileWithRetry: ファイルを削除しつつ、失敗したらリトライする
func removeFileWithRetry(filePath string) {
	for i := range MAX_RETRIES {
		err := os.Remove(filePath)
		if err == nil {
			log.Println("ファイルを削除:", filePath)
			return
		}

		log.Printf("ファイルの削除に失敗 (%d/%d): %v", i+1, 5, err)
		time.Sleep(3 * time.Second)
	}

	log.Println("リトライ上限に達したため、ファイルの削除に失敗:", filePath)
}

// moveFileWithRetry: ファイルを移動しつつ、失敗したらリトライする
func moveFileWithRetry(sourcePath, destPath string) {
	for i := range MAX_RETRIES {
		err := os.Rename(sourcePath, destPath)
		if err == nil {
			log.Println("ファイルを移動:", destPath)
			return
		}

		log.Printf("ファイルの移動に失敗 (%d/%d): %v", i+1, MAX_RETRIES, err)
		time.Sleep(RETRY_DELAY)
	}

	log.Println("リトライ上限に達したため、ファイルの移動に失敗:", destPath)
}

// createWorldNameDir: ワールド名のディレクトリを作成
func createWorldNameDir(targetDirPath, worldName string) error {
	// ワールド名をディレクトリ名に使用可能な形に変換
	safeWorldName := convertToSafeDirectoryName(worldName)
	worldDirPath := filepath.Join(targetDirPath, safeWorldName)

	// TargetDirPath以下にワールド名のディレクトリを作成
	if _, err := os.Stat(worldDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(worldDirPath, os.ModePerm); err != nil {
			return fmt.Errorf("ワールド名のディレクトリの作成に失敗: %w", err)
		}
	}

	log.Println("ワールド名のディレクトリを作成:", worldDirPath)

	return nil
}
