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

// createWorldNameDir: ワールド名のディレクトリを作成
func createWorldNameDir(targetDirPath, worldName string) error {
	safeWorldName := convertToSafeDirectoryName(worldName)
	worldDirPath := filepath.Join(targetDirPath, safeWorldName)

	if _, err := os.Stat(worldDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(worldDirPath, os.ModePerm); err != nil {
			return fmt.Errorf("ワールド名のディレクトリの作成に失敗: %w", err)
		}
	}

	log.Println("ワールド名のディレクトリを作成:", worldDirPath)

	return nil
}

// openFileWithRetry: ファイルを開きつつ、失敗したらリトライする
func openFileWithRetry(filePath string) (*os.File, error) {
	for i := 0; i < MAX_RETRIES; i++ {
		file, err := os.Open(filePath)
		if err == nil {
			log.Println("ファイルを開きました:", filePath)
			return file, nil
		}

		log.Printf("ファイルのオープンに失敗。リトライします (%d/%d): %v", i+1, MAX_RETRIES, err)
		time.Sleep(RETRY_DELAY)
	}

	return nil, fmt.Errorf("リトライ上限に達しました: %s", filePath)
}

// removeFileWithRetry: ファイルを削除しつつ、失敗したらリトライする
func removeFileWithRetry(filePath string) {
	for i := range MAX_RETRIES {
		err := os.Remove(filePath)
		if err == nil {
			log.Println("ファイルを削除:", filePath)
			return
		}

		log.Printf("ファイルの削除に失敗。リトライします (%d/%d): %v", i+1, 5, err)
		time.Sleep(3 * time.Second)
	}

	log.Println("リトライ上限に達しました:", filePath)
}

// moveFileWithRetry: ファイルを移動しつつ、失敗したらリトライする
func moveFileWithRetry(sourcePath, destPath string) {
	for i := range MAX_RETRIES {
		err := os.Rename(sourcePath, destPath)
		if err == nil {
			log.Println("ファイルを移動:", destPath)
			return
		}

		log.Printf("ファイルの移動に失敗。リトライします (%d/%d): %v", i+1, MAX_RETRIES, err)
		time.Sleep(RETRY_DELAY)
	}

	log.Println("リトライ上限に達しました:", destPath)
}
