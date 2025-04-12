package app

import (
	"log"

	"fyne.io/systray"
	"github.com/pkg/browser"
	"github.com/sqweek/dialog"
)

// setupMenu: メニューの設定
func (a *App) setupMenu() {
	mSettings := systray.AddMenuItem("設定", "設定を変更する")
	mVRCLogDir := mSettings.AddSubMenuItem("ログフォルダを指定", "VRChatのログフォルダを指定する")
	mVRCPhotoDir := mSettings.AddSubMenuItem("写真フォルダを指定", "VRChatの写真フォルダを指定する")
	mAbout := systray.AddMenuItem("About", "アプリについて")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("終了", "アプリを終了する")

	// イベントをハンドリング
	go func() {
		for {
			select {
			case <-mAbout.ClickedCh:
				browser.OpenURL("https://github.com/arrow2nd/memento")
			case <-mVRCLogDir.ClickedCh:
				a.UpdateVRCLogDir()
			case <-mVRCPhotoDir.ClickedCh:
				if a.UpdateVRCPictureDir() {
					(*a.watcher).Setup()
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

// UpdateVRCLogDir: VRChatのログフォルダを選択して更新する
func (a *App) UpdateVRCLogDir() {
	dir, err := dialog.Directory().SetStartDir(a.config.VRCLogDirPath).Title("VRChatのログフォルダを指定").Browse()
	if err != nil {
		log.Println("ログフォルダの選択に失敗:", err)
		return
	}

	if err := a.config.SetVRCLogDirPath(dir); err != nil {
		log.Println("ログフォルダの設定に失敗:", err)
		return
	}

	log.Println("ログフォルダの設定を更新:", dir)
}

// UpdateVRCPictureDir: VRChatの写真フォルダを選択して更新する
func (a *App) UpdateVRCPictureDir() bool {
	dir, err := dialog.Directory().SetStartDir(a.config.PictureDirPath).Title("VRChatの写真フォルダを指定").Browse()
	if err != nil {
		log.Println("写真フォルダの選択に失敗:", err)
		return false
	}

	if err := a.config.SetRootDirPath(dir); err != nil {
		log.Println("写真フォルダの設定に失敗:", err)
		return false
	}

	log.Println("写真フォルダの設定を更新:", dir)

	return true
}
