package picture

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// createWorldNameDir: ワールド名のディレクトリを作成
func createWorldNameDir(targetDirPath, worldName string) error {
	// ワールド名をディレクトリ名に使用可能な形に変換
	safeWorldName := convertToSafeDirectoryName(worldName)
	worldDirPath := filepath.Join(targetDirPath, safeWorldName)

	// TargetDirPath以下にワールド名のディレクトリを作成
	if _, err := os.Stat(worldDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(worldDirPath, 0755); err != nil {
			return fmt.Errorf("ワールド名のディレクトリを作成できませんでした: %w", err)
		}
	}

	return nil
}

// moveFile: ファイルを移動（コピー＆削除）
func moveFile(srcPath, destPath string) error {
	// ソースファイルを開く
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("ファイルを開けませんでした: %w", err)
	}
	defer src.Close()

	// 移動先ファイルを作成
	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("移動先ファイルを作成できませんでした: %w", err)
	}
	defer dst.Close()

	// ファイルの内容をコピー
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("ファイルのコピーに失敗しました: %w", err)
	}

	// 元のファイルを削除
	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("元のファイルを削除できませんでした: %w", err)
	}

	return nil
}

