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
	mVRCLogDir := mSettings.AddSubMenuItem("ログフォルダを指定", "VRChatのログフォルダを指定します")
	mVRCPhotoDir := mSettings.AddSubMenuItem("写真フォルダを指定", "VRChatの写真フォルダを指定します")
	mConvertToJpeg := mSettings.AddSubMenuItemCheckbox("JPEGに変換", "写真をJPEGに変換します。ついでに撮影日時なども書き込みます", a.config.ConvertToJpeg)
	mAbout := systray.AddMenuItem("アプリについて", "アプリの配布ページを開きます")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("終了", "アプリを終了する")

	toggleConvertToJpeg := func() {
		if a.config.ConvertToJpeg {
			mConvertToJpeg.Check()
		} else {
			mConvertToJpeg.Uncheck()
		}
	}

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
			case <-mConvertToJpeg.ClickedCh:
				a.UpdateConvertToJpeg()
				toggleConvertToJpeg()

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

// UpdateConvertToJpeg: JPEG変換の設定を更新する
func (a *App) UpdateConvertToJpeg() {
	if err := a.config.SetConvertToJpeg(!a.config.ConvertToJpeg); err != nil {
		log.Println("JPEG変換の設定に失敗:", err)
		return
	}

	log.Println("JPEG変換の設定を更新:", a.config.ConvertToJpeg)

}
