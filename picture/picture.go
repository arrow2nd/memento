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
func MoveToWorldNameDir(opts MoveToWorldNameDirOpts, convertToJpeg bool, jpegQuality int) error {
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

	// 設定が有効かつ、マルチレイヤーの写真でない場合はJPEGに変換
	if convertToJpeg && !isMultiLayerPicture(opts.PicturePath) {
		if _, err := encodeJpegWithExif(opts.PicturePath, worldDirPath, opts.WorldVisit.Name, takePictureTime, jpegQuality); err != nil {
			return fmt.Errorf("ファイルの移動に失敗: %w", err)
		}

		go removeFileWithRetry(opts.PicturePath)

		return nil
	}

	// PNG画像をそのまま移動
	destPath := filepath.Join(worldDirPath, filepath.Base(opts.PicturePath))
	go moveFileWithRetry(opts.PicturePath, destPath)

	return nil
}
