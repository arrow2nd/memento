package logparser

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"time"
)

// WorldVisit: ワールドの訪問履歴
type WorldVisit struct {
	// Time: 訪問した日時
	Time time.Time
	// Name: 訪問したワールド名
	Name string
}

// findLatestWorldVisitFromLog: ログファイルから最も直近に訪問したワールドを探す
func findLatestWorldVisitFromLog(logPath string) (*WorldVisit, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 逆順に読むためにファイルサイズを取得
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	maxLineSize := int64(4096)
	var latestVisit *WorldVisit

	// 例: 2025.04.09 21:16:28 Debug      -  [Behaviour] Entering Room: 夜と、共に。
	re := regexp.MustCompile(`^(\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2}).*Entering Room: (.+)$`)

	// 後ろから読む
	for pos := fileSize; pos > 0 && latestVisit == nil; {
		readSize := maxLineSize
		if pos < readSize {
			readSize = pos
		}

		startPos := pos - readSize
		buffer := make([]byte, readSize)

		_, err := file.ReadAt(buffer, startPos)
		if err != nil {
			return nil, err
		}

		lines := bufio.NewScanner(bufio.NewReader(bytes.NewReader(buffer)))

		for lines.Scan() {
			line := lines.Text()
			matches := re.FindStringSubmatch(line)

			if len(matches) == 3 {
				timestampStr := matches[1]
				worldName := matches[2]

				// ローカルタイムゾーンとして日時をパース
				// NOTE: VRCのログの日付もローカルっぽいので
				t, err := time.ParseInLocation("2006.01.02 15:04:05", timestampStr, time.Local)
				if err != nil {
					continue
				}

				latestVisit = &WorldVisit{
					Time: t,
					Name: worldName,
				}
				break
			}
		}

		pos = startPos
	}

	if latestVisit == nil {
		return nil, errors.New("ワールド訪問履歴が見つかりませんでした")
	}

	return latestVisit, nil
}
