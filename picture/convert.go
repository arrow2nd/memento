package picture

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// encodeJpeg: JPEGにエンコードして保存する
func encodeJpeg(srcPath, destDirPath string) error {
	// 入力元から出力先のパスを生成
	fileName := filepath.Base(srcPath)
	newFileName := strings.Replace(fileName, filepath.Ext(fileName), ".jpg", 1)
	jpgPath := filepath.Join(destDirPath, newFileName)

	// 元画像を開く
	file, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("元画像へのアクセス失敗: %w", err)
	}
	defer file.Close()

	log.Println("元画像を開きました:", srcPath)

	// PNGをデコード
	img, err := png.Decode(file)
	if err != nil {
		return fmt.Errorf("PNGのデコードに失敗: %w", err)
	}

	log.Println("画像をデコードしました:", srcPath)

	// 出力先ファイルを作成
	out, err := os.Create(jpgPath)
	if err != nil {
		return fmt.Errorf("出力先ファイルの作成に失敗: %w", err)
	}
	defer out.Close()

	log.Println("出力先ファイルを作成:", jpgPath)

	// JPEGとしてエンコード
	options := jpeg.Options{
		Quality: 90,
	}
	if err := jpeg.Encode(out, img, &options); err != nil {
		return fmt.Errorf("JPEG画像の変換に失敗: %w", err)
	}

	log.Println("JPEG画像に変換:", jpgPath)

	return nil
}

