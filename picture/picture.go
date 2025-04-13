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
		log.Printf("ファイルの移動に失敗: %v, バックグラウンドでリトライします", err)

		go retryMoveFile(opts.PicturePath, destPath)

		return nil
	}

	log.Println("ファイルを移動:", destPath)
	return nil
}
