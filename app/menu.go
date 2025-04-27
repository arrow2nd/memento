package app

import (
	"log"

	"fyne.io/systray"
	"github.com/pkg/browser"
	"github.com/sqweek/dialog"
)

// setupMenu: メニューの設定
func (a *App) setupMenu() {
	// 設定サブメニュー
	mSettings := systray.AddMenuItem("設定", "設定を変更する")

	// 機能設定
	mConvertToJpeg := mSettings.AddSubMenuItemCheckbox("JPEGに変換", "写真をJPEGに変換します。ついでに撮影日時なども書き込みます", a.config.ConvertToJpeg)
	mSettings.AddSeparator()

	// フォルダ設定
	mVRCLogDir := mSettings.AddSubMenuItem("ログフォルダを指定", "VRChatのログフォルダを指定します")
	mVRCPhotoDir := mSettings.AddSubMenuItem("写真フォルダを指定", "VRChatの写真フォルダを指定します")
	mSettings.AddSeparator()

	// mementoの設定ディレクトリを開く
	mOpenConfigDir := mSettings.AddSubMenuItem("設定フォルダを開く", "mementoの設定フォルダを開きます")

	// その他
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
			case <-mOpenConfigDir.ClickedCh:
				browser.OpenURL(a.config.ConfigDirPath)

			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

// UpdateVRCLogDir: VRChatのログフォルダを選択して更新する
func (a *App) UpdateVRCLogDir() {
	title := "VRChatのログフォルダを指定してください。\n\"output_log_なんとか.txt\" みたいなファイルが置いてあるフォルダです。"

	dir, err := dialog.Directory().SetStartDir(a.config.VRCLogDirPath).Title(title).Browse()
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
	title := "VRChatの写真フォルダを指定してください。\n\"2025-04\" みたいなフォルダではなく、その1つ上のフォルダです。(たぶん \"VRChat\" って名前のはず)"

	dir, err := dialog.Directory().SetStartDir(a.config.PictureDirPath).Title(title).Browse()
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
