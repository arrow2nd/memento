package app

import (
	"log"

	"fyne.io/systray"
	"github.com/arrow2nd/memento/config"
	"github.com/arrow2nd/memento/watcher"
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
	cfg := config.New()

	return &App{
		name:    name,
		version: version,
		config:  cfg,
		watcher: watcher.New(cfg),
	}
}

func (a *App) Run() {
	systray.Run(a.onReady, a.onExit)
}

func (a *App) onReady() {
	systray.SetTitle(a.name)
	systray.SetTooltip("VRCの写真フォルダを監視中です")

	// TODO: アイコンの埋め込みもしたい

	// メニューアイテムの設定
	mQuit := systray.AddMenuItem("終了", "アプリを終了する")

	// 監視を開始
	go a.watcher.Start()

	// 終了イベントをハンドリング
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func (a *App) onExit() {
	log.Println("終了しています")
}
