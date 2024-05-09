package window

import (
	"log"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	internal "github.com/HackJack14/SteamSync/Client/Internal"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type ClientWindowEvent func(*ClientWindow)

type ClientWindow struct {
	OkClicked     ClientWindowEvent
	CancelClicked ClientWindowEvent
	OnIPReceived  ClientWindowEvent
	app           fyne.App
	window        fyne.Window
	ipadress      string
	upTab         UploadTab
	downTab       DownloadTab
	contTabs      *container.AppTabs
	btnOk         *widget.Button
	btnCancel     *widget.Button
}

type UploadTab struct {
	labelNewGame       *widget.Label
	entryNewGame       *widget.Entry
	btnSearchGameSaves *widget.Button
	labelExistingGame  *widget.Label
	selectExistingGame *widget.Select
	labelDirectory     *widget.Label
	entryDirectory     *widget.Entry
	btnOpenDialog      *widget.Button
}

type DownloadTab struct {
	labelGame      *widget.Label
	selectGame     *widget.Select
	labelDirectory *widget.Label
	entryDirectory *widget.Entry
	btnOpenDialog  *widget.Button
}

func NewClientWindow() *ClientWindow {
	return &ClientWindow{}
}

func (cla *ClientWindow) Init() {
	cla.app = app.New()
	cla.app.Settings().SetTheme(theme.DarkTheme())
	setScaling()
	cla.window = cla.app.NewWindow("SteamSync Client")
	cla.window.Resize(fyne.NewSize(500, 200))
	cla.window.SetMaster()

	cla.downTab.labelGame = widget.NewLabel("Game:")
	cla.downTab.selectGame = widget.NewSelect([]string{}, func(s string) {})
	cla.downTab.labelDirectory = widget.NewLabel("Save Directory")
	cla.downTab.entryDirectory = widget.NewEntry()
	cla.downTab.btnOpenDialog = widget.NewButton("...", func() { cla.openDirDialog(cla.downTab.entryDirectory) })

	cla.upTab.labelNewGame = widget.NewLabel("New Game:")
	cla.upTab.entryNewGame = widget.NewEntry()
	cla.upTab.btnSearchGameSaves = widget.NewButton("Search", func() { cla.searchGameSaves(cla.upTab.entryNewGame.Text, cla.upTab.entryDirectory) })
	cla.upTab.labelExistingGame = widget.NewLabel("Existing Game")
	cla.upTab.selectExistingGame = widget.NewSelect([]string{}, func(s string) {})
	cla.upTab.labelDirectory = widget.NewLabel("Save Directory")
	cla.upTab.entryDirectory = widget.NewEntry()
	cla.upTab.btnOpenDialog = widget.NewButton("...", func() { cla.openDirDialog(cla.upTab.entryDirectory) })

	cla.btnOk = widget.NewButton("Ok", func() { cla.OkClicked(cla) })
	cla.btnCancel = widget.NewButton("Cancel", func() { cla.CancelClicked(cla) })

	contBtns := container.NewHBox(cla.btnOk, cla.btnCancel)

	contNewGame := container.NewStack(cla.upTab.entryNewGame, container.NewHBox(layout.NewSpacer(), cla.upTab.btnSearchGameSaves))
	contFolderOpenUpload := container.NewStack(cla.upTab.entryDirectory, container.NewHBox(layout.NewSpacer(), cla.upTab.btnOpenDialog))
	contUpload := container.NewVBox(cla.upTab.labelNewGame, contNewGame, cla.upTab.labelExistingGame, cla.upTab.selectExistingGame, cla.upTab.labelDirectory, contFolderOpenUpload)

	contFolderOpenDownload := container.NewStack(cla.downTab.entryDirectory, container.NewHBox(layout.NewSpacer(), cla.downTab.btnOpenDialog))
	contDownload := container.NewVBox(cla.downTab.labelGame, cla.downTab.selectGame, cla.downTab.labelDirectory, contFolderOpenDownload)

	cla.contTabs = container.NewAppTabs(container.NewTabItem("Upload", contUpload), container.NewTabItem("Download", contDownload))
	contMain := container.NewVBox(cla.contTabs, contBtns)
	cla.window.SetContent(contMain)
}

func (cla *ClientWindow) Show() {
	entryIP := widget.NewEntry()
	dialog.ShowForm("IPAddress", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Input the servers ipaddress", entryIP)}, func(b bool) {
		if b {
			cla.ipadress = entryIP.Text
			cla.OnIPReceived(cla)
		} else {
			cla.Close()
		}
	}, cla.window)
	cla.window.ShowAndRun()
}

func (cla *ClientWindow) Close() {
	cla.window.Close()
	cla.app.Quit()
}

func (cla *ClientWindow) openDirDialog(outputEntry *widget.Entry) {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			log.Fatalln(err)
		}

		if uri != nil {
			outputEntry.SetText(uri.Path())
		}
	}, cla.window)
}

func (cla *ClientWindow) searchGameSaves(game string, output *widget.Entry) {
	windows, linux, err := internal.GetSaves(game)
	if err != nil {
		log.Println(err)
	}

	save := ""
	switch runtime.GOOS {
	case "windows":
		save = windows
	case "linux":
		save = linux
	}

	output.SetText(internal.ExractFullPath(save))
}

func setScaling() {
	var actScreenWidth int

	err := glfw.Init()
	if err != nil {
		actScreenWidth = 100
	}
	defer glfw.Terminate()

	monitor := glfw.GetPrimaryMonitor()
	actScreenWidth, _ = monitor.GetPhysicalSize()

	log.Println(actScreenWidth)

	if actScreenWidth <= 100 {
		os.Setenv("FYNE_SCALE", "0.2")
	} else if actScreenWidth > 100 && actScreenWidth <= 200 {
		os.Setenv("FYNE_SCALE", "0.75")
	} else {
		os.Setenv("FYNE_SCALE", "1")
	}
}

func (cla *ClientWindow) SetGames(games []string) {
	cla.upTab.selectExistingGame.Options = games
	cla.downTab.selectGame.Options = games
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
		if cla.upTab.entryNewGame.Text != "" {
			return cla.upTab.entryNewGame.Text
		} else {
			return cla.upTab.selectExistingGame.Selected
		}
	} else {
		return cla.downTab.selectGame.Selected
	}
}

func (cla *ClientWindow) GetDir() string {
	if cla.contTabs.Selected().Text == "Upload" {
		return cla.upTab.entryDirectory.Text
	} else {
		return cla.downTab.entryDirectory.Text
	}
}
