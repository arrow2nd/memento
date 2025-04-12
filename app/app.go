package app

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/arrow2nd/memento/config"
	"github.com/arrow2nd/memento/watcher"
)

// App: アプリケーション
type App struct {
	name    string
	version string
	config  *config.Config
	window  *fyne.Window
	watcher *watcher.Watcher
}

// New: 作成
func New(name, version string) *App {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("configの初期化に失敗: ", err)
	}

	return &App{
		name:    name,
		version: version,
		config:  cfg,
		window:  nil,
		watcher: watcher.New(cfg),
	}
}

// Run: 実行
func (a *App) Run() {
	app := app.New()
	a.window = a.configureWindow(&app)

	// デスクトップアプリケーションとして動作しているか確認
	if desk, ok := app.(desktop.App); ok {
		a.configureSystemTray(desk)
	}

	// 監視を開始
	go a.watcher.Start()

	app.Run()
}
