package picture

import (
	"fmt"
	"image"
	"image/png"
	"log"
)

// decodePng: PNG画像をデコード
func decodePng(srcPath string) (*image.Image, error) {
	file, err := openFileWithRetry(srcPath)
	if err != nil {
		return nil, fmt.Errorf("PNG画像のオープンに失敗: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("PNG画像のデコードに失敗: %w", err)
	}

	log.Println("PNG画像をデコード:", srcPath)

	return &img, nil
}
