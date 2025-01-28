package window

import (
	"log"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/HackJack14/SteamSync/network"
)

var (
	OnSelectedPeer func(peer network.Peer)
	OnUploadGame   func(game string) error
	OnDownloadGame func(game string) error
)

type peerButton struct {
	peer   network.Peer
	button widget.Clickable
}

type gameButton struct {
	game   string
	button widget.Clickable
}

type step = int

const (
	peer step = iota
	game
	sync
)

type peerSelector struct {
	buttonList []peerButton
	list       layout.List
}

type gameSelector struct {
	buttonsList []gameButton
	list        layout.List
}

type syncSelector struct {
	flex     layout.Flex
	upload   widget.Clickable
	download widget.Clickable
}

func RenderWindow() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func readPeers(peers chan []network.Peer, stop chan bool, window *app.Window) {
	for {
		stopping := false
		select {
		case stopping = <-stop:
		default:
			stopping = false
		}
		if stopping {
			return
		}

		newPeers, err := network.GetAllPeers()
		if err != nil {
			return
		}
		peers <- newPeers
		window.Invalidate()
	}
}

func readGames() ([]string, error) {
	entries, err := os.ReadDir(filepath.Dir(os.Args[0]))
	if err != nil {
		return []string{}, err
	}

	games := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			games = append(games, entry.Name())
		}
	}
	return games, nil
}

var (
	stepCount    step   = 0
	selectedGame string = ""
	pSelector           = peerSelector{
		list: layout.List{
			Axis:        layout.Vertical,
			Alignment:   layout.Middle,
			ScrollToEnd: false,
		},
	}
	gSelector = gameSelector{
		list: layout.List{
			Axis:        layout.Vertical,
			Alignment:   layout.Middle,
			ScrollToEnd: false,
		},
	}
	sSelector = syncSelector{
		flex: layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceSides,
		},
	}
)

func run(window *app.Window) error {
	theme := material.NewTheme()
	stopDiscovery := make(chan bool)
	peerChan := make(chan []network.Peer, 1)
	go readPeers(peerChan, stopDiscovery, window)
	gamesList, err := readGames()
	gSelector.buttonsList = make([]gameButton, 0)
	for _, game := range gamesList {
		button := gameButton{
			game:   game,
			button: widget.Clickable{},
		}
		gSelector.buttonsList = append(gSelector.buttonsList, button)
	}
	if err != nil {
		return err
	}
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			switch stepCount {
			case peer:
				var peers []network.Peer
				select {
				case peers = <-peerChan:
					pSelector.buttonList = make([]peerButton, len(peers))
					for i, p := range peers {
						pSelector.buttonList[i].peer = p
					}
				default:
				}
				pSelector.Layout(gtx, theme)
			case game:
				gSelector.Layout(gtx, theme)
			case sync:
				sSelector.Layout(gtx, theme)
			}

			e.Frame(gtx.Ops)
		}
	}
}

func (selector *peerSelector) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return selector.list.Layout(gtx, len(pSelector.buttonList), func(gtx layout.Context, index int) layout.Dimensions {
		if selector.buttonList[index].Clicked(gtx) {
			stepCount = game
		}
		return selector.buttonList[index].Layout(gtx, theme)
	})
}

func (button *peerButton) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return material.Button(theme, &button.button, button.peer.Hostname).Layout(gtx)
}

func (button *peerButton) Clicked(gtx layout.Context) bool {
	return button.button.Clicked(gtx)
}

func (button *gameButton) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return material.Button(theme, &button.button, button.game).Layout(gtx)
}

func (button *gameButton) Clicked(gtx layout.Context) bool {
	return button.button.Clicked(gtx)
}

func (selector *gameSelector) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return selector.list.Layout(gtx, len(selector.buttonsList), func(gtx layout.Context, index int) layout.Dimensions {
		if selector.buttonsList[index].Clicked(gtx) {
			stepCount = sync
			selectedGame = gSelector.buttonsList[index].game
		}
		return selector.buttonsList[index].Layout(gtx, theme)
	})
}

func (selector *syncSelector) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Spacing:   layout.SpaceSides,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if selector.upload.Clicked(gtx) {
				OnUploadGame(selectedGame)
			}
			return material.Button(theme, &selector.upload, "Upload").Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if selector.download.Clicked(gtx) {
				OnDownloadGame(selectedGame)
			}
			return material.Button(theme, &selector.download, "Downlaod").Layout(gtx)
		}),
	)
}
