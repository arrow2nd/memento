package picture

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arrow2nd/memento/logparser"
)

type MoveToWorldNameDirOpts struct {
	PicturePath   string
	TargetDirPath string
	WorldVisit    *logparser.WorldVisit
}

// MoveToWorldNameDir: 写真をワールド名のディレクトリに移動
func MoveToWorldNameDir(opts MoveToWorldNameDirOpts) error {
	// 撮影日時を取得
	takePictureTime, err := getPictureSaveDate(opts.PicturePath)
	if err != nil {
		return fmt.Errorf("撮影日時の取得に失敗: %w", err)
	}

	// 撮影日時がワールド訪問日時よりも前なら中断
	if takePictureTime.Before(opts.WorldVisit.Time) {
		return fmt.Errorf("撮影日時がワールド訪問日時以前のためスキップ: %s", opts.PicturePath)
	}

	// ワールド名のディレクトリを作成
	if err := createWorldNameDir(opts.TargetDirPath, opts.WorldVisit.Name); err != nil {
		return err
	}

	// 移動先のパスを生成
	safeWorldName := convertToSafeDirectoryName(opts.WorldVisit.Name)
	worldDirPath := filepath.Join(opts.TargetDirPath, safeWorldName)
	pictureName := filepath.Base(opts.PicturePath)
	destPath := filepath.Join(worldDirPath, pictureName)

	// ファイルを移動
	if err := os.Rename(opts.PicturePath, destPath); err != nil {
		return fmt.Errorf("ファイルの移動に失敗: %w", err)
	}

	log.Println("ファイルを移動:", destPath)
	return nil
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
