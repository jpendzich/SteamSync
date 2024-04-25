package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type BtnFunc func()

var OkClicked BtnFunc
var CancelClicked BtnFunc

var EntryIP *widget.Entry
var EntryRequest *widget.Entry
var EntryGame *widget.Entry
var EntryDir *widget.Entry

var a fyne.App
var w fyne.Window

func CreateWindow() {
	a = app.New()
	w = a.NewWindow("SteamSync Client")
	w.Resize(fyne.NewSize(500, 200))

	labelIP := widget.NewLabel("IP-Address")
	EntryIP = widget.NewEntry()
	labelRequest := widget.NewLabel("Request")
	EntryRequest = widget.NewEntry()
	labelGame := widget.NewLabel("Game")
	EntryGame = widget.NewEntry()
	labelDir := widget.NewLabel("Directory")
	EntryDir = widget.NewEntry()
	btnOk := widget.NewButton("Ok", OkClicked)
	btnCancel := widget.NewButton("Cancel", CancelClicked)
	contBtns := container.NewHBox(btnOk, btnCancel)
	contMain := container.NewVBox(labelIP, EntryIP, labelRequest, EntryRequest, labelGame, EntryGame, labelDir, EntryDir, contBtns)
	w.SetContent(contMain)
	w.ShowAndRun()
}

func CloseWindow() {
	w.Hide()
	a.Quit()
}
