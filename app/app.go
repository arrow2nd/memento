package app

import (
	_ "embed"
	"fmt"
	"log"

	"fyne.io/systray"
	"github.com/arrow2nd/memento/config"
	"github.com/arrow2nd/memento/logger"
	"github.com/arrow2nd/memento/watcher"
	"github.com/gofrs/flock"
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
	name     string
	version  string
	config   *config.Config
	watcher  *watcher.Watcher
	lockFile *flock.Flock
}

func New() *App {
	// ロガーを初期化
	log.SetOutput(logger.Setup(appName))

	// 設定を初期化
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

	// 重複起動防止
	app.checkAlreadyRunning()

	// 設定されたディレクトリを確認
	app.checkDirectories()

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
	if a.lockFile != nil {
		if err := a.lockFile.Unlock(); err != nil {
			log.Println("ロックファイルの解放に失敗:", err)
		}
	}

	log.Println("アプリケーションを終了")
}
