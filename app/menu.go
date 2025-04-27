package app

import (
	"fyne.io/systray"
	"github.com/arrow2nd/memento/autostart"
	"github.com/pkg/browser"
)

// setupMenu: メニューの設定
func (a *App) setupMenu() {
	// 設定サブメニュー
	mSettings := systray.AddMenuItem("設定", "設定を変更する")

	// 機能設定
	mAutoStart := mSettings.AddSubMenuItemCheckbox(
		"自動起動",
		"Windowsの起動時に自動起動します",
		autostart.IsAutoStartEnabled(a.name),
	)
	mConvertToJpeg := mSettings.AddSubMenuItemCheckbox(
		"JPEGに変換",
		"写真をJPEGに変換します。ついでに撮影日時なども書き込みます",
		a.config.ConvertToJpeg,
	)
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

	// イベントをハンドリング
	go func() {
		for {
			select {
			case <-mAbout.ClickedCh:
				browser.OpenURL("https://github.com/arrow2nd/memento")

			case <-mVRCLogDir.ClickedCh:
				a.updateVRCLogDir()

			case <-mVRCPhotoDir.ClickedCh:
				if a.updateVRCPictureDir() {
					(*a.watcher).Setup()
				}

			case <-mConvertToJpeg.ClickedCh:
				a.UpdateConvertToJpeg(mConvertToJpeg)

			case <-mOpenConfigDir.ClickedCh:
				browser.OpenURL(a.config.ConfigDirPath)

			case <-mAutoStart.ClickedCh:
				a.toggleAutoStart(mAutoStart)

			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}
