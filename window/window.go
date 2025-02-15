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
	"github.com/HackJack14/SteamSync/internal"
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

type backButton struct {
	inset  layout.Inset
	button widget.Clickable
}

type gameSearchBar struct {
	inset      layout.Inset
	lastSearch string
	edit       widget.Editor
}

type gameScrollBar struct {
	scrollbar widget.Scrollbar
	start     float32
	end       float32
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
	buttonsList   []gameButton
	list          layout.List
	scrollbar     gameScrollBar
	searchBar     gameSearchBar
	searchFlex    layout.Flex
	scrollbarFlex layout.Flex
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
	access       internal.DbAccess
	stepCount    step        = game
	selectedGame string      = ""
	frame        layout.Flex = layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}
	backbutton backButton = backButton{
		inset: layout.Inset{
			Top: 5,
		},
		button: widget.Clickable{},
	}
	pSelector = peerSelector{
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
		scrollbar: gameScrollBar{},
		searchBar: gameSearchBar{
			lastSearch: "",
			inset: layout.Inset{
				Top:    5,
				Bottom: 5,
				Left:   5,
				Right:  5,
			},
			edit: widget.Editor{
				SingleLine: true,
				Submit:     true,
			},
		},
		searchFlex: layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEnd,
		},
		scrollbarFlex: layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceBetween,
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
	access = *internal.NewDbAccess()
	defer access.Close()
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

			frame.Layout(
				gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					var dimensions layout.Dimensions
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
						dimensions = pSelector.Layout(gtx, theme)
					case game:
						dimensions = gSelector.Layout(gtx, theme)
					case sync:
						dimensions = sSelector.Layout(gtx, theme)
					}
					return dimensions
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if backbutton.Clicked(gtx) {
						switch stepCount {
						case peer:
							stepCount = sync
						case game:
							stepCount = game
						case sync:
							stepCount = game
						}
					}
					return backbutton.Layout(gtx, theme)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func (selector *peerSelector) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return selector.list.Layout(gtx, len(pSelector.buttonList), func(gtx layout.Context, index int) layout.Dimensions {
		if selector.buttonList[index].Clicked(gtx) {
			OnSelectedPeer(selector.buttonList[index].peer)
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

func (button *backButton) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return button.inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Button(theme, &button.button, "back").Layout(gtx)
	})
}

func (button *backButton) Clicked(gtx layout.Context) bool {
	return button.button.Clicked(gtx)
}

func (searchBar *gameSearchBar) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return searchBar.inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Editor(theme, &searchBar.edit, "search for a game").Layout(gtx)
	})
}

func (scrollBar *gameScrollBar) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return material.Scrollbar(theme, &scrollBar.scrollbar).Layout(gtx, layout.Vertical, scrollBar.start, scrollBar.end)
}

func (selector *gameSelector) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	// only set lastSearch if edit.Text() changed
	if selector.searchBar.lastSearch != selector.searchBar.edit.Text() {
		selector.searchBar.lastSearch = selector.searchBar.edit.Text()
		games := access.GetGameMatchingPattern(selector.searchBar.lastSearch)
		selector.buttonsList = make([]gameButton, len(games))
		for i, game := range games {
			selector.buttonsList[i] = gameButton{
				game: game,
			}
		}
	}
	// check for clicked button
	for i := 0; i < len(selector.buttonsList); i++ {
		if selector.buttonsList[i].Clicked(gtx) {
			stepCount = sync
			selectedGame = selector.buttonsList[i].game
		}
	}

	if selector.scrollbar.scrollbar.Dragging() {
		selector.list.ScrollBy(float32(len(selector.buttonsList)) * selector.scrollbar.scrollbar.ScrollDistance())
	}

	if len(selector.buttonsList) > 0 {
		selector.scrollbar.start = float32(selector.list.Position.First) / float32(len(selector.buttonsList))
		selector.scrollbar.end = selector.scrollbar.start + float32(selector.list.Position.Count)/float32(len(selector.buttonsList))
	} else {
		selector.scrollbar.start = 0
		selector.scrollbar.end = 1
	}

	return selector.searchFlex.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return selector.searchBar.Layout(gtx, theme)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return selector.scrollbarFlex.Layout(
				gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return selector.list.Layout(gtx, len(selector.buttonsList), func(gtx layout.Context, index int) layout.Dimensions {
						return selector.buttonsList[index].Layout(gtx, theme)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return selector.scrollbar.Layout(gtx, theme)
				}),
			)

		}),
	)
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
