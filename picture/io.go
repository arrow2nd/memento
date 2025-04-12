package picture

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

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

// isFileLocked: ファイルがロックされているかどうかを確認
func isFileLocked(filePath string) bool {
	// 書き込みモードで開く
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0666)

	// エラーが発生した場合、ファイルがロックされているとみなす
	if err != nil {
		return true
	}

	file.Close()
	return false
}

// moveFile: ファイルを移動（コピー＆削除）
func moveFile(srcPath, destPath string) error {
	maxRetries := 3
	retryDelay := 2 * time.Second

	for i := range maxRetries {
		if isFileLocked(srcPath) {
			if i == maxRetries-1 {
				return fmt.Errorf("ファイルが他のプロセスによってロックされています: %s", srcPath)
			}

			// リトライ前に少し待機
			time.Sleep(retryDelay)

			continue
		}

		// ファイルを移動
		if err := os.Rename(srcPath, destPath); err != nil {
			return fmt.Errorf("ファイルの移動に失敗しました: %w", err)
		}

		log.Println("ファイルを移動しました: ", destPath)

		return nil
	}

	return fmt.Errorf("最大リトライ回数を超えました")
}
