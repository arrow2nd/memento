package app

import (
	"fmt"
	"log"

	"fyne.io/systray"
	"github.com/arrow2nd/memento/config"
	"github.com/arrow2nd/memento/watcher"
	"github.com/sqweek/dialog"
)

// App: アプリケーション
type App struct {
	name    string
	version string
	config  *config.Config
	watcher *watcher.Watcher
}

// New: 作成
func New(name, version string) *App {
	cfg, err := config.New(name)
	if err != nil {
		dialog.Message("設定を取得できませんでした").Title("エラー").Error()
		log.Fatal("configの初期化に失敗: ", err)
	}

	app := &App{
		name:    name,
		version: version,
		config:  cfg,
	}

	// 設定されたディレクトリが存在するか確認
	pictureDirExists, logDirExists := app.config.CheckDirectoriesExist()

	if !pictureDirExists {
		dialog.Message("%s\n%s\n\n%s\n%s",
			"写真フォルダが見つかりませんでした。",
			"次の画面でVRChatの写真フォルダを選んでください。",
			"※選ぶのは写真が直接入っているフォルダではなく、その親フォルダです。",
			"（写真が見える場所より1つ上のフォルダを選んでください）",
		).Title("写真フォルダの確認").Info()

		app.UpdateVRCPictureDir()
	}

	if !logDirExists {
		dialog.Message("%s\n%s",
			"VRChatのログフォルダが見つかりませんでした。",
			"次の画面でVRChatのログフォルダを選んでください。",
		).Title("ログフォルダの確認").Info()

		app.UpdateVRCLogDir()
	}

	// 監視処理の初期化
	watcher, err := watcher.New(cfg)
	if err != nil {
		dialog.Message("監視を開始できませんでした").Title("エラー").Error()
		log.Fatal("watcherの初期化に失敗: ", err)
	}

	app.watcher = watcher

	return app
}

func (a *App) Run() {
	systray.Run(a.onReady, a.onExit)
}

func (a *App) onReady() {
	systray.SetTitle(a.name)
	a.updateTooltip()

	// TODO: アイコンの埋め込みもしたい

	// メニューの設定
	a.setupMenu()

	// 監視を開始
	go a.watcher.Start()
}

func (a *App) onExit() {
	log.Println("終了しています")
}

// updateTooltip: ツールチップを更新
func (a *App) updateTooltip() {
	tooltip := fmt.Sprint(
		fmt.Sprintf("%s v.%s\n", a.name, a.version),
		fmt.Sprintf("写真フォルダ: %s\n", a.config.PictureDirPath),
		fmt.Sprintf("ログフォルダ: %s", a.config.VRCLogDirPath),
	)

	systray.SetTooltip(tooltip)
}
