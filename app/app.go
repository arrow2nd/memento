package app

import (
	"fmt"
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
	cfg, err := config.New(name)
	if err != nil {
		log.Fatal("configの初期化に失敗: ", err)
	}

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
