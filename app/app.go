package app

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	lockFile   *flock.Flock
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

	// 重複起動防止のためのロックファイルを設定
	lockFilePath := filepath.Join(cfg.ConfigDirPath, appName+".lock")
	lockFile = flock.New(lockFilePath)

	// 重複起動を防止
	if app.isAlreadyRunning() {
		os.Exit(0)
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
	if lockFile != nil {
		err := lockFile.Unlock()
		if err != nil {
			log.Println("ロックファイルの解放に失敗:", err)
		}
	}

	log.Println("アプリケーションを終了")
}

// isAlreadyRunning: 重複起動チェック
func (a *App) isAlreadyRunning() bool {
	locked, err := lockFile.TryLock()
	if err != nil {
		log.Println("ロックファイルの作成に失敗:", err)
		return false
	}

	if !locked {
		dialog.Message("%sは既に起動しています！\nタスクトレイを確認してみてください。", a.name).Title("起動エラー").Error()
		log.Println("既に起動しているため終了")
		return true
	}

	return false
}
