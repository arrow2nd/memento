package picture

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arrow2nd/memento/logparser"
)

// findNearestWorldVisit: 写真の保存日時に最も近いワールド訪問を取得
func findNearestWorldVisit(pictureTime time.Time, visits []logparser.WorldVisit) (*logparser.WorldVisit, error) {
	if len(visits) == 0 {
		return nil, errors.New("ワールド訪問履歴がありません")
	}

	var nearestVisit *logparser.WorldVisit
	var minTimeDiff time.Duration = time.Duration(math.MaxInt64)

	for i, visit := range visits {
		var timeDiff time.Duration

		if pictureTime.After(visit.Time) {
			timeDiff = pictureTime.Sub(visit.Time)
		} else {
			timeDiff = visit.Time.Sub(pictureTime)
		}

		if timeDiff < minTimeDiff {
			minTimeDiff = timeDiff
			nearestVisit = &visits[i]
		}
	}

	if nearestVisit == nil {
		return nil, errors.New("適切なワールド訪問が見つかりませんでした")
	}
	return nearestVisit, nil
}


// getPictureSaveDate: 写真の保存日時を取得
func getPictureSaveDate(path string) (time.Time, error) {
	// ファイル名から日時の抽出してみる
	fileName := filepath.Base(path)
	timeFromName, err := extractDateFromFileName(fileName)
	if err == nil {
		return timeFromName, nil
	}

	// ファイル名からの抽出に失敗した場合は、ファイルの修正日時を使用
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}

// extractDateFromFileName: ファイル名から日時を抽出
func extractDateFromFileName(fileName string) (time.Time, error) {
	// VRChat_YYYY-MM-DD_HH-MM-SS.ms_WIDTHxHEIGHT 形式を想定
	parts := strings.Split(fileName, "_")
	if len(parts) < 3 || !strings.HasPrefix(parts[0], "VRChat") {
		return time.Time{}, fmt.Errorf("サポートされていないファイル名形式です: %s", fileName)
	}

	datePart := parts[1] // YYYY-MM-DD
	timePart := parts[2] // HH-MM-SS.ms

	// YYYY-MM-DD HH:MM:SS.ms 形式に
	timePartFormatted := strings.Replace(timePart, "-", ":", 2)
	timeStr := fmt.Sprintf("%s %s", datePart, timePartFormatted)

	// 日時をパース (2006-01-02 15:04:05.999)
	t, err := time.ParseInLocation("2006-01-02 15:04:05.999", timeStr, time.Local)
	if err != nil {
		// 一応ミリ秒なしも試す
		t, err = time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
		if err != nil {
			return time.Time{}, fmt.Errorf("ファイル名から日時を解析できませんでした: %w", err)
		}
	}

	return t, nil
}
