package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ClientWindowEvent func(*ClientWindow)

type ClientWindow struct {
	OkClicked     ClientWindowEvent
	CancelClicked ClientWindowEvent
	GamesClicked  ClientWindowEvent
	app           fyne.App
	window        fyne.Window
	ipwindow      ipAddressWindow
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
	btnGames      *widget.Button
}

type ipAddressWindow struct {
	window    fyne.Window
	labelIP   *widget.Label
	entryIP   *widget.Entry
	btnOk     *widget.Button
	btnCancel *widget.Button
}

func NewClientWindow() *ClientWindow {
	return &ClientWindow{}
}

func newIPAddressWindow() *ipAddressWindow {
	return &ipAddressWindow{}
}

func (cla *ClientWindow) Init() {
	cla.app = app.New()
	cla.window = cla.app.NewWindow("SteamSync Client")
	cla.window.Resize(fyne.NewSize(500, 200))
	cla.window.SetMaster()
	cla.ipwindow = *newIPAddressWindow()
	cla.ipwindow.window = cla.app.NewWindow("IPAddress")
	cla.ipwindow.labelIP = widget.NewLabel("IPAddress:")
	cla.ipwindow.entryIP = widget.NewEntry()
	cla.ipwindow.btnOk = widget.NewButton("Ok", func() { cla.ipadress = cla.ipwindow.entryIP.Text })
	cla.ipwindow.btnCancel = widget.NewButton("Cancel", func() {
		cla.ipwindow.window.Close()
		cla.window.Close()
		cla.app.Quit()
	})
	contIPBtns := container.NewHBox(cla.ipwindow.btnOk, cla.ipwindow.btnCancel)
	contIP := container.NewVBox(cla.ipwindow.labelIP, cla.ipwindow.entryIP, contIPBtns)
	cla.ipwindow.window.SetContent(contIP)

	cla.labelRequest = widget.NewLabel("Choose what to do")
	cla.selectRequest = widget.NewSelect([]string{"UPLOAD", "DOWNLOAD"}, func(s string) {})
	cla.labelGame = widget.NewLabel("Game:")
	cla.selectGames = widget.NewSelect([]string{}, func(s string) {})
	cla.entryGame = widget.NewEntry()
	cla.labelDir = widget.NewLabel("Directory:")
	cla.entryDir = widget.NewEntry()
	cla.btnOk = widget.NewButton("Ok", func() { cla.OkClicked(cla) })
	cla.btnCancel = widget.NewButton("Cancel", func() { cla.CancelClicked(cla) })
	cla.btnGames = widget.NewButton("Update Games", func() { cla.GamesClicked(cla) })

	contBtns := container.NewHBox(cla.btnOk, cla.btnCancel, cla.btnGames)
	contMain := container.NewVBox(cla.labelRequest, cla.selectRequest,
		cla.labelGame, cla.selectGames, cla.entryGame, cla.labelDir, cla.entryDir, contBtns)
	cla.window.SetContent(contMain)
}

func (cla *ClientWindow) Show() {
	cla.ipwindow.window.Show()
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
