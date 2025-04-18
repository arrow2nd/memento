package app

import (
	_ "embed"
	"fmt"
	"log"

	"fyne.io/systray"
	"github.com/arrow2nd/memento/config"
	"github.com/arrow2nd/memento/logger"
	"github.com/arrow2nd/memento/watcher"
	"github.com/sqweek/dialog"
)

var (
	appName    = "memento"
	appVersion = "develop"
)

//go:embed trayicon.ico
var trayIcon []byte

// App: アプリケーション
type App struct {
	name    string
	version string
	config  *config.Config
	watcher *watcher.Watcher
}

func New() *App {
	// ロガーを初期化
	log.SetOutput(logger.Setup(appName))

	cfg, err := config.New(appName)
	if err != nil {
		dialog.Message("設定を取得できませんでした").Title("エラー").Error()
		log.Fatal("configの初期化に失敗:", err)
	}

	app := &App{
		name:    appName,
		version: appVersion,
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
		log.Fatal("watcherの初期化に失敗:", err)
	}

	app.watcher = watcher

	return app
}

func (a *App) Run() {
	log.Println("起動:", a.version)

	systray.Run(a.onReady, a.onExit)
}

func (a *App) onReady() {
	systray.SetIcon(trayIcon)
	systray.SetTitle(a.name)
	systray.SetTooltip(fmt.Sprintf("%s %s", a.name, a.version))

	a.setupMenu()

	go a.watcher.Start()
}

func (a *App) onExit() {
	log.Println("アプリケーションを終了")
}
