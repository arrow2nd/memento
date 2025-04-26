package picture

import (
	"fmt"
	"time"

	"github.com/dsoprea/go-exif/v3"
)

// setExifData: 共通のEXIFデータを設定
func setExifData(rootBuilder *exif.IfdBuilder, worldName string, pictureTime time.Time) error {
	ifdBuilder, err := exif.GetOrCreateIbFromRootIb(rootBuilder, "IFD")
	if err != nil {
		return fmt.Errorf("ifdBuilderの作成に失敗: %w", err)
	}

	// ワールド名をタイトルとして設定
	err = ifdBuilder.SetStandardWithName("ImageDescription", fmt.Sprintf("撮影ワールド: %s", worldName))
	if err != nil {
		return fmt.Errorf("ImageDescriptionの設定に失敗: %w", err)
	}

	exifBuilder, err := exif.GetOrCreateIbFromRootIb(rootBuilder, "IFD/Exif")
	if err != nil {
		return fmt.Errorf("exifBuilderの作成に失敗: %w", err)
	}

	// 撮影日時を設定
	dateTimeStr := pictureTime.Format("2006:01:02 15:04:05")
	err = exifBuilder.SetStandardWithName("DateTimeOriginal", dateTimeStr)
	if err != nil {
		return fmt.Errorf("DateTimeOriginalの設定に失敗: %w", err)
	}

	return nil
}
