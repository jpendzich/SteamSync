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
	upTab         *UploadTab
	downTab       *DownloadTab
	contTabs      *container.AppTabs
	btnOk         *widget.Button
	btnCancel     *widget.Button
}

type UploadTab struct {
	labelNewGame       *widget.Label
	entryNewGame       *widget.Entry
	labelExistingGame  *widget.Label
	selectExistingGame *widget.Select
	labelDirectory     *widget.Label
	entryDirectory     *widget.Entry
}

type DownloadTab struct {
	labelGame      *widget.Label
	selectGame     *widget.Select
	labelDirectory *widget.Label
	entryDirectory *widget.Entry
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

	cla.downTab.labelGame = widget.NewLabel("Game:")
	cla.downTab.selectGame = widget.NewSelect([]string{}, func(s string) {})
	cla.downTab.labelDirectory = widget.NewLabel("Save Directory")
	cla.downTab.entryDirectory = widget.NewEntry()

	cla.upTab.labelNewGame = widget.NewLabel("New Game:")
	cla.upTab.entryNewGame = widget.NewEntry()
	cla.upTab.labelExistingGame = widget.NewLabel("Existing Game")
	cla.upTab.selectExistingGame = widget.NewSelect([]string{}, func(s string) {})
	cla.upTab.labelDirectory = widget.NewLabel("Save Directory")
	cla.upTab.entryDirectory = widget.NewEntry()

	cla.btnOk = widget.NewButton("Ok", func() { cla.OkClicked(cla) })
	cla.btnCancel = widget.NewButton("Cancel", func() { cla.CancelClicked(cla) })

	contBtns := container.NewHBox(cla.btnOk, cla.btnCancel)
	contUpload := container.NewVBox(cla.upTab.labelNewGame, cla.upTab.entryNewGame, cla.upTab.labelExistingGame, cla.upTab.selectExistingGame, cla.upTab.labelDirectory, cla.upTab.entryDirectory)
	contDownload := container.NewVBox(cla.downTab.labelGame, cla.downTab.selectGame, cla.downTab.labelDirectory, cla.downTab.entryDirectory)
	cla.contTabs = container.NewAppTabs(container.NewTabItem("Upload", contUpload), container.NewTabItem("Download", contDownload))
	contMain := container.NewVBox(cla.contTabs, contBtns)
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
	if cla.contTabs.Selected().Text == "Upload" {
		cla.upTab.selectExistingGame.Options = games
	} else {
		cla.downTab.selectGame.Options = games
	}
}

func (cla *ClientWindow) GetIP() string {
	return cla.ipadress
}

func (cla *ClientWindow) GetRequest() string {
	if cla.contTabs.Selected().Text == "Upload" {
		return "UPLOAD"
	} else {
		return "DOWNLOAD"
	}
}

func (cla *ClientWindow) GetGame() string {
	if cla.contTabs.Selected().Text == "Upload" {
		return cla.upTab.entry http.ResponseWriter, r *http.Request
	} else {
		cla.downTab.selectGame.Options = games
	}
	return cla.entryGame.Text
}

func (cla *ClientWindow) GetDir() string {
	return cla.entryDir.Text
}
