package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/sqweek/dialog"
)

// checkDirectories : 設定されたディレクトリを確認し、必要に応じてユーザーに選択を促す
func (a *App) checkDirectories() {
	// 設定されたディレクトリが存在するか確認
	pictureDirExists, logDirExists := a.config.CheckDirectoriesExist()

	if !pictureDirExists {
		dialog.Message("%s\n%s\n\n%s\n%s",
			"写真フォルダが見つかりませんでした。",
			"次の画面でVRChatの写真フォルダを選んでください。",
			"※選ぶのは写真が直接入っているフォルダではなく、その親フォルダです。",
			"（写真が見える場所より1つ上のフォルダを選んでください）",
		).Title("写真フォルダの確認").Info()

		a.updateVRCPictureDir()
	}

	if !logDirExists {
		dialog.Message("%s\n%s",
			"VRChatのログフォルダが見つかりませんでした。",
			"次の画面でVRChatのログフォルダを選んでください。",
		).Title("ログフォルダの確認").Info()

		a.updateVRCLogDir()
	}
}

// checkAlreadyRunning: 既に起動しているか確認
func (a *App) checkAlreadyRunning() {
	// ロックファイルを作成
	if a.lockFile == nil {
		lockFilePath := filepath.Join(a.config.ConfigDirPath, appName+".lock")
		a.lockFile = flock.New(lockFilePath)
	}

	locked, err := a.lockFile.TryLock()
	if err != nil {
		log.Println("ロックファイルの作成に失敗:", err)
		return // とりあえず起動させとく
	}

	// ロックされていたら終了
	if !locked {
		dialog.Message("%sは既に起動しています！\nタスクトレイを確認してみてください。", a.name).Title("起動エラー").Error()
		log.Println("既に起動しているため終了")
		os.Exit(0)
	}
}

