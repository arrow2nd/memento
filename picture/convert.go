package picture

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
)

// encodeJpegWithExif: JPEGにエンコードしてEXIFデータも同時に書き込む
func encodeJpegWithExif(srcPath, destDirPath string, worldName string, pictureTime time.Time, quality int) (string, error) {
	// 入力元から出力先のパスを生成
	fileName := filepath.Base(srcPath)
	newFileName := strings.Replace(fileName, filepath.Ext(fileName), ".jpg", 1)
	jpgPath := filepath.Join(destDirPath, newFileName)

	// 元画像を開く
	file, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("元画像へのアクセス失敗: %w", err)
	}
	defer file.Close()

	log.Println("元画像を開きました:", srcPath)

	// PNGをデコード
	img, err := png.Decode(file)
	if err != nil {
		return "", fmt.Errorf("PNGのデコードに失敗: %w", err)
	}

	log.Println("画像をデコードしました:", srcPath)

	// 一時バッファにJPEGとしてエンコード
	options := jpeg.Options{
		Quality: quality,
	}

	// 一旦メモリ上にJPEGをエンコード
	jpegBuf := new(strings.Builder)
	if err := jpeg.Encode(jpegBuf, img, &options); err != nil {
		return "", fmt.Errorf("JPEG画像の変換に失敗: %w", err)
	}

	log.Println("JPEG画像に変換しました")

	// JPEGデータをバイト配列に変換
	jpegData := []byte(jpegBuf.String())

	// パーサーを作成
	jmp := jpegstructure.NewJpegMediaParser()

	// JPEGデータを解析してセグメントリストを得る
	ec, err := jmp.ParseBytes(jpegData)
	if err != nil {
		return "", fmt.Errorf("JPEG解析に失敗: %w", err)
	}
	sl := ec.(*jpegstructure.SegmentList)

	// IfdBuilderを作成
	rootBuilder, err := sl.ConstructExifBuilder()
	if err != nil {
		return "", fmt.Errorf("EXIFビルダーの作成に失敗: %w", err)
	}

	// 基本的な EXIF 情報を設定 (IFD)
	ifdBuilder, err := exif.GetOrCreateIbFromRootIb(rootBuilder, "IFD")
	if err != nil {
		return "", fmt.Errorf("IFD ビルダーの作成に失敗: %w", err)
	}

	// ワールド名をタイトルとして設定
	err = ifdBuilder.SetStandardWithName("ImageDescription", fmt.Sprintf("撮影ワールド: %s", worldName))
	if err != nil {
		return "", fmt.Errorf("ImageDescription の設定に失敗: %w", err)
	}

	// EXIF サブディレクトリを取得
	exifBuilder, err := exif.GetOrCreateIbFromRootIb(rootBuilder, "IFD/Exif")
	if err != nil {
		return "", fmt.Errorf("EXIF ビルダーの作成に失敗: %w", err)
	}

	// 撮影日時を設定
	dateTimeStr := pictureTime.Format("2006:01:02 15:04:05")
	err = exifBuilder.SetStandardWithName("DateTimeOriginal", dateTimeStr)
	if err != nil {
		return "", fmt.Errorf("DateTimeOriginal の設定に失敗: %w", err)
	}

	// SegmentList を更新
	err = sl.SetExif(rootBuilder)
	if err != nil {
		return "", fmt.Errorf("EXIF データの設定に失敗: %w", err)
	}

	// 出力先ファイルを作成
	out, err := os.Create(jpgPath)
	if err != nil {
		return "", fmt.Errorf("出力先ファイルの作成に失敗: %w", err)
	}
	defer out.Close()

	log.Println("出力先ファイルを作成:", jpgPath)

	// EXIFデータを含むJPEG画像を書き込み
	err = sl.Write(out)
	if err != nil {
		return "", fmt.Errorf("EXIFデータ付きJPEG画像の書き込みに失敗: %w", err)
	}

	log.Println("EXIF付きJPEG画像を保存:", jpgPath)

	return jpgPath, nil
}
