package logparser

import (
	"bufio"
	"os"
	"regexp"
	"time"
)

type WorldVisit struct {
	// Time: 訪問した日時
	Time time.Time
	// WorldName: 訪問したワールド名
	WorldName string
}

// parseWorldVisitsFromLog: ログファイルからワールド訪問履歴を解析
func parseWorldVisitsFromLog(logPath string) ([]WorldVisit, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var visits []WorldVisit
	scanner := bufio.NewScanner(file)

	// 例: 2025.04.09 21:16:28 Debug      -  [Behaviour] Entering Room: 夜と、共に。
	re := regexp.MustCompile(`^(\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2}).*Entering Room: (.+)$`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)

		if len(matches) == 3 {
			timestampStr := matches[1]
			worldName := matches[2]

			// ローカルタイムゾーンとして日時を解析
			t, err := time.ParseInLocation("2006.01.02 15:04:05", timestampStr, time.Local)
			if err != nil {
				continue
			}

			visits = append(visits, WorldVisit{
				Time:      t,
				WorldName: worldName,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}
