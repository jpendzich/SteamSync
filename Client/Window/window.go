package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ClientWindowEvent func(*ClientWindow)

type ClientWindow struct {
	OkClicked     ClientWindowEvent
	CancelClicked ClientWindowEvent
	OnIPReceived  ClientWindowEvent
	app           fyne.App
	window        fyne.Window
	ipadress      string
	labelRequest  *widget.Label
	selectRequest *widget.Select
	labelGame     *widget.Label
	selectGames   *widget.Select
	entryGame     *widget.Entry
	labelDir      *widget.Label
	entryDir      *widget.Entry
	btnOk         *widget.Button
	btnCancel     *widget.Button
}

func NewClientWindow() *ClientWindow {
	return &ClientWindow{}
}

func (cla *ClientWindow) Init() {
	cla.app = app.New()
	cla.app.Settings().SetTheme(theme.DarkTheme())
	cla.window = cla.app.NewWindow("SteamSync Client")
	cla.window.Resize(fyne.NewSize(500, 200))
	cla.window.SetMaster()

	cla.labelRequest = widget.NewLabel("Choose what to do")
	cla.selectRequest = widget.NewSelect([]string{"UPLOAD", "DOWNLOAD"}, func(s string) {})
	cla.labelGame = widget.NewLabel("Game:")
	cla.selectGames = widget.NewSelect([]string{}, func(s string) {})
	cla.entryGame = widget.NewEntry()
	cla.labelDir = widget.NewLabel("Directory:")
	cla.entryDir = widget.NewEntry()
	cla.btnOk = widget.NewButton("Ok", func() { cla.OkClicked(cla) })
	cla.btnCancel = widget.NewButton("Cancel", func() { cla.CancelClicked(cla) })

	contBtns := container.NewHBox(cla.btnOk, cla.btnCancel)
	contMain := container.NewVBox(cla.labelRequest, cla.selectRequest,
		cla.labelGame, cla.selectGames, cla.entryGame, cla.labelDir, cla.entryDir, contBtns)
	cla.window.SetContent(contMain)
}

func (cla *ClientWindow) Show() {
	entryIP := widget.NewEntry()
	dialog.ShowForm("IPAddress", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Input the servers ipaddress", entryIP)}, func(b bool) {
		if b {
			cla.ipadress = entryIP.Text
		} else {
			cla.Close()
		}
	}, cla.window)
	cla.OnIPReceived(cla)
	cla.window.ShowAndRun()
}

func (cla *ClientWindow) Close() {
	cla.window.Close()
	cla.app.Quit()
}

func (cla *ClientWindow) SetGames(games []string) {
	cla.selectGames.Options = games
}

func (cla *ClientWindow) GetIP() string {
	return cla.ipadress
}

func (cla *ClientWindow) GetRequest() string {
	return cla.selectRequest.Selected
}

func (cla *ClientWindow) GetGame() string {
	return cla.entryGame.Text
}

func (cla *ClientWindow) GetDir() string {
	return cla.entryDir.Text
}
