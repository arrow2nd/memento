package app

import (
	"log"

	"fyne.io/systray"
	"github.com/arrow2nd/memento/autostart"
	"github.com/sqweek/dialog"
)

// updateVRCLogDir: VRChatのログフォルダを選択して更新する
func (a *App) updateVRCLogDir() {
	title := "VRChatのログフォルダを指定してください。\n\"output_log_なんとか.txt\" みたいなファイルが置いてあるフォルダです。"

	dir, err := dialog.Directory().SetStartDir(a.config.VRCLogDirPath).Title(title).Browse()
	if err != nil {
		log.Println("ログフォルダの選択に失敗:", err)
		return
	}

	if err := a.config.SetVRCLogDirPath(dir); err != nil {
		dialog.Message("ログフォルダの設定に失敗しました").Title("エラー").Error()
		log.Println("ログフォルダの設定に失敗:", err)
		return
	}

	log.Println("ログフォルダの設定を更新:", dir)
}

// updateVRCPictureDir: VRChatの写真フォルダを選択して更新する
func (a *App) updateVRCPictureDir() bool {
	title := "VRChatの写真フォルダを指定してください。\n\"2025-04\" みたいなフォルダではなく、その1つ上のフォルダです。(たぶん \"VRChat\" って名前のはず)"

	dir, err := dialog.Directory().SetStartDir(a.config.PictureDirPath).Title(title).Browse()
	if err != nil {
		log.Println("写真フォルダの選択に失敗:", err)
		return false
	}

	if err := a.config.SetRootDirPath(dir); err != nil {
		dialog.Message("写真フォルダの設定に失敗しました").Title("エラー").Error()
		log.Println("写真フォルダの設定に失敗:", err)
		return false
	}

	log.Println("写真フォルダの設定を更新:", dir)

	return true
}

// UpdateConvertToJpeg: JPEG変換の設定を更新する
func (a *App) UpdateConvertToJpeg(mConvertToJpeg *systray.MenuItem) {
	if err := a.config.SetConvertToJpeg(!a.config.ConvertToJpeg); err != nil {
		dialog.Message("JPEG変換の設定に失敗しました").Title("エラー").Error()
		log.Println("JPEG変換の設定に失敗:", err)
		return
	}

	// メニューの表示を更新
	if a.config.ConvertToJpeg {
		mConvertToJpeg.Check()
	} else {
		mConvertToJpeg.Uncheck()
	}

	log.Println("JPEG変換の設定を更新:", a.config.ConvertToJpeg)

}

// toggleAutoStart: 自動起動の設定を切り替える
func (a *App) toggleAutoStart(menuItem *systray.MenuItem) {
	newSetting := !autostart.IsAutoStartEnabled(a.name)

	// レジストリを更新
	err := autostart.SetAutoStart(a.name, newSetting)
	if err != nil {
		dialog.Message("自動起動の設定に失敗しました").Title("エラー").Error()
		log.Println("自動起動の設定に失敗:", err)
		return
	}

	// メニューの表示を更新
	if newSetting {
		menuItem.Check()
	} else {
		menuItem.Uncheck()
	}

	log.Println("自動起動の設定を変更:", newSetting)
}

