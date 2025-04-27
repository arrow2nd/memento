package picture

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dsoprea/go-jpeg-image-structure/v2"
)

// encodeJpegWithExif: JPEGに変換しつつEXIFを追加
func encodeJpegWithExif(srcPath, destDirPath string, worldName string, pictureTime time.Time, quality int) (string, error) {
	jpgPath := generateJpegPath(srcPath, destDirPath)

	// PNGをデコード
	img, err := decodePng(srcPath)
	if err != nil {
		return "", err
	}

	// JPEGに変換
	jpegData, err := encodeToJpegBytes(*img, quality)
	if err != nil {
		return "", err
	}

	// EXIFを追加
	sl, err := addExifData(jpegData, worldName, pictureTime)
	if err != nil {
		return "", err
	}

	// 書き込み
	if err := writeJpegWithExif(sl, jpgPath); err != nil {
		return "", err
	}

	return jpgPath, nil
}

// generateJpegPath: JPEG画像の出力パスを生成
func generateJpegPath(srcPath, destDirPath string) string {
	fileName := filepath.Base(srcPath)
	newFileName := strings.Replace(fileName, filepath.Ext(fileName), ".jpg", 1)

	return filepath.Join(destDirPath, newFileName)
}

// encodeToJpegBytes: 画像をJPEGのバイト配列に変換
func encodeToJpegBytes(img image.Image, quality int) ([]byte, error) {
	options := jpeg.Options{
		Quality: quality,
	}

	jpegBuf := new(strings.Builder)
	if err := jpeg.Encode(jpegBuf, img, &options); err != nil {
		return nil, fmt.Errorf("JPEGへの変換に失敗: %w", err)
	}

	log.Printf("JPEGへ変換 (Quality: %d)\n", quality)

	return []byte(jpegBuf.String()), nil
}

// addExifData: EXIFをJPEGに追加
func addExifData(jpegData []byte, worldName string, pictureTime time.Time) (*jpegstructure.SegmentList, error) {
	jmp := jpegstructure.NewJpegMediaParser()

	// セグメントリスト取得
	ec, err := jmp.ParseBytes(jpegData)
	if err != nil {
		return nil, fmt.Errorf("JPEGの解析に失敗: %w", err)
	}
	sl := ec.(*jpegstructure.SegmentList)

	rootBuilder, err := sl.ConstructExifBuilder()
	if err != nil {
		return nil, fmt.Errorf("rootBuilderの作成に失敗: %w", err)
	}

	// 共通EXIF設定処理を呼び出し
	if err := setExifData(rootBuilder, worldName, pictureTime); err != nil {
		return nil, err
	}

	// SegmentListを更新
	err = sl.SetExif(rootBuilder)
	if err != nil {
		return nil, fmt.Errorf("EXIFの設定に失敗: %w", err)
	}

	return sl, nil
}

// writeJpegWithExif: EXIF付きでJPEG画像を保存
func writeJpegWithExif(sl *jpegstructure.SegmentList, jpgPath string) error {
	out, err := os.Create(jpgPath)
	if err != nil {
		return fmt.Errorf("出力先ファイルの作成に失敗: %w", err)
	}
	defer out.Close()

	log.Println("出力先ファイルを作成:", jpgPath)

	err = sl.Write(out)
	if err != nil {
		return fmt.Errorf("EXIF付きのJPEGの書き込みに失敗: %w", err)
	}

	log.Println("EXIF付きのJPEGを保存:", jpgPath)

	return nil
}
