package window

import (
	"log"
	"os"
	"path"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type gameButton struct {
	game   string
	button widget.Clickable
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

func readGames() ([]string, error) {
	entries, err := os.ReadDir(path.Dir(os.Args[0]))
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
	gamesList   []string
	buttonsList []gameButton
)

func run(window *app.Window) error {
	theme := material.NewTheme()
	gamesList, err := readGames()
	buttonsList = make([]gameButton, 0)
	for _, game := range gamesList {
		button := gameButton{
			game:   game,
			button: widget.Clickable{},
		}
		buttonsList = append(buttonsList, button)
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

			renderFirst(gtx, theme)

			e.Frame(gtx.Ops)
		}
	}
}

func renderFirst(gtx layout.Context, theme *material.Theme) {
	list := layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: false,
		Alignment:   layout.Middle,
	}
	list.Layout(gtx, len(gamesList), func(gtx layout.Context, index int) layout.Dimensions {
		return material.Button(theme, &buttonsList[index].button, buttonsList[index].game).Layout(gtx)
	})
}

func renderSecond(gtx layout.Context, theme *material.Theme, game string) {
	// downloadButton := new(widget.Clickable)
	// uploadButton := new(widget.Clickable)
}
