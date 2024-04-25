package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ClientAppEvent func(*ClientApp)

type ClientApp struct {
	OkClicked     ClientAppEvent
	CancelClicked ClientAppEvent
	app           fyne.App
	window        fyne.Window
	labelIP       *widget.Label
	entryIP       *widget.Entry
	labelRequest  *widget.Label
	selectRequest *widget.Select
	labelGame     *widget.Label
	entryGame     *widget.Entry
	labelDir      *widget.Label
	entryDir      *widget.Entry
	btnOk         *widget.Button
	btnCancel     *widget.Button
}

func NewClientApp() *ClientApp {
	return &ClientApp{}
}

func (cla *ClientApp) Init() {
	cla.app = app.New()
	cla.window = cla.app.NewWindow("SteamSync Client")
	cla.window.Resize(fyne.NewSize(500, 200))

	cla.labelIP = widget.NewLabel("IP-Address:")
	cla.entryIP = widget.NewEntry()
	cla.labelRequest = widget.NewLabel("Choose what to do")
	cla.selectRequest = widget.NewSelect([]string{"UPLOAD", "DOWNLOAD"}, func(s string) {})
	cla.labelGame = widget.NewLabel("Game:")
	cla.entryGame = widget.NewEntry()
	cla.labelDir = widget.NewLabel("Directory:")
	cla.entryDir = widget.NewEntry()
	cla.btnOk = widget.NewButton("Ok", func() { cla.OkClicked(cla) })
	cla.btnCancel = widget.NewButton("Cancel", func() { cla.CancelClicked(cla) })

	contBtns := container.NewHBox(cla.btnOk, cla.btnCancel)
	contMain := container.NewVBox(cla.labelIP, cla.entryIP, cla.labelRequest, cla.selectRequest,
		cla.labelGame, cla.entryGame, cla.labelDir, cla.entryDir, contBtns)
	cla.window.SetContent(contMain)
}

func (cla *ClientApp) Show() {
	cla.window.ShowAndRun()
}

func (cla *ClientApp) Close() {
	cla.window.Close()
	cla.app.Quit()
}

func (cla *ClientApp) GetIP() string {
	return cla.entryIP.Text
}

func (cla *ClientApp) GetRequest() string {
	return cla.selectRequest.Selected
}

func (cla *ClientApp) GetGame() string {
	return cla.entryGame.Text
}

func (cla *ClientApp) GetDir() string {
	return cla.entryDir.Text
}
