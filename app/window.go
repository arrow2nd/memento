package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// configureWindow: メインウィンドウの設定
func (a *App) configureWindow(app *fyne.App) *fyne.Window {
	window := (*app).NewWindow(a.name)

	window.SetContent(widget.NewLabel("VRCの写真フォルダを監視中です"))

	// ウィンドウを閉じる代わりに非表示に
	window.SetCloseIntercept(func() {
		window.Hide()
	})

	return &window
}

// configureSystemTray: システムトレイメニューの設定
func (a *App) configureSystemTray(desk desktop.App) {
	menu := fyne.NewMenu(a.name,
		fyne.NewMenuItem("設定", func() {
			a.showConfigWindow()
		}),
	)

	desk.SetSystemTrayMenu(menu)
}

// showConfigWindow: 設定ウィンドウを表示
func (a *App) showConfigWindow() {
	(*a.window).Show()
}

