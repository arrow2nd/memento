package picture

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/arrow2nd/memento/logparser"
)

type MoveToWorldNameDirOpts struct {
	PicturePath   string
	TargetDirPath string
	WorldVisit    *logparser.WorldVisit
}

// MoveToWorldNameDir: 写真をワールド名のディレクトリに移動
func MoveToWorldNameDir(opts MoveToWorldNameDirOpts, convertToJpeg bool) error {
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

	// PNG画像ではない場合はスキップ
	if filepath.Ext(opts.PicturePath) != ".png" {
		log.Println("PNG画像ではないためスキップ:", opts.PicturePath)
		return nil
	}

	// 移動先のパスを生成
	safeWorldName := convertToSafeDirectoryName(opts.WorldVisit.Name)
	worldDirPath := filepath.Join(opts.TargetDirPath, safeWorldName)

	// JPEGに変換する場合
	if convertToJpeg {
		if err := encodeJpeg(opts.PicturePath, worldDirPath); err != nil {
			return err
		}

		// 移動元の画像を削除
		go removeFileWithRetry(opts.PicturePath)

		return nil
	}

	// 変換しないなら、そのまま移動する
	pictureName := filepath.Base(opts.PicturePath)
	destPath := filepath.Join(worldDirPath, pictureName)
	go moveFileWithRetry(opts.PicturePath, destPath)

	return nil
}
