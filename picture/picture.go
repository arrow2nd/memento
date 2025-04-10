package picture

import (
	"fmt"
	"path/filepath"

	"github.com/arrow2nd/memento/logparser"
)

type MoveToWorldNameDirOpts struct {
	PicturePath   string
	TargetDirPath string
	WorldVisits   []logparser.WorldVisit
}

// MoveToWorldNameDir: 写真をワールド名のディレクトリに移動
func MoveToWorldNameDir(opts MoveToWorldNameDirOpts) error {
	// 写真の保存日時を取得
	saveTime, err := getPictureSaveDate(opts.PicturePath)
	if err != nil {
		return fmt.Errorf("写真の保存日時を取得できませんでした: %w", err)
	}

	// ワールド訪問履歴と照合して、保存日時に最も近いワールド訪問を取得
	visit, err := findNearestWorldVisit(saveTime, opts.WorldVisits)
	if err != nil {
		return fmt.Errorf("写真に最も近いワールド訪問を見つけられませんでした: %w", err)
	}

	// ワールド名のディレクトリを作成
	if err := createWorldNameDir(opts.TargetDirPath, visit.WorldName); err != nil {
		return err
	}

	// 移動先のパスを生成
	safeWorldName := convertToSafeDirectoryName(visit.WorldName)
	worldDirPath := filepath.Join(opts.TargetDirPath, safeWorldName)
	pictureName := filepath.Base(opts.PicturePath)
	destPath := filepath.Join(worldDirPath, pictureName)

	// ファイルを移動
	if err := moveFile(opts.PicturePath, destPath); err != nil {
		return err
	}

	return nil
}

