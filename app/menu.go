package app

import (
	"log"

	"fyne.io/systray"
	"github.com/sqweek/dialog"
)

// setupMenu: メニューの設定
func (a *App) setupMenu() {
	mSettings := systray.AddMenuItem("設定", "設定を変更する")
	mVRCLogDir := mSettings.AddSubMenuItem("ログフォルダを指定", "VRChatのログフォルダを指定する")
	mVRCPhotoDir := mSettings.AddSubMenuItem("写真フォルダを指定", "VRChatの写真フォルダを指定する")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("終了", "アプリを終了する")

	// イベントをハンドリング
	go func() {
		for {
			select {
			case <-mVRCLogDir.ClickedCh:
				a.onVRCLogDirClick()
			case <-mVRCPhotoDir.ClickedCh:
				a.onVRCPhotoDirClick()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

// onVRCLogDirClick: VRChatのログフォルダを指定する
func (a *App) onVRCLogDirClick() {
	dir, err := dialog.Directory().SetStartDir(a.config.VRCLogDirPath).Title("VRChatのログフォルダを指定").Browse()
	if err != nil {
		log.Println("ログフォルダの選択に失敗:", err)
		return
	}

	if err := a.config.SetVRCLogDirPath(dir); err != nil {
		log.Println("ログフォルダの設定に失敗:", err)
		return
	}

	log.Println("ログフォルダを設定しました:", dir)
	a.updateTooltip()
}

// onVRCPhotoDirClick: VRChatの写真フォルダを指定する
func (a *App) onVRCPhotoDirClick() {
	dir, err := dialog.Directory().SetStartDir(a.config.PictureDirPath).Title("VRChatの写真フォルダを指定").Browse()
	if err != nil {
		log.Println("写真フォルダの選択に失敗:", err)
		return
	}

	if err := a.config.SetRootDirPath(dir); err != nil {
		log.Println("写真フォルダの設定に失敗:", err)
		return
	}

	log.Println("写真フォルダを設定しました:", dir)
	a.updateTooltip()
}
