package ui

import (
	"erik-schuetze/ypsilon-14/pkg/util"
	"fmt"

	"github.com/jroimartin/gocui"
)

var initialValues = map[string]bool{
	"legendView":    false,
	"operationView": false,
}
var SubViews = util.NewMutexMap(initialValues)

var MenuItems = []string{"Station Overview", "Docking Bay", "Floorplan", "Base Operations", "Scan Admin ID Card"}

func InitLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if titleView, err := g.SetView("title", 2, 1, maxX-7, 3); err != nil {
		titleView.Title = " MINING STATION TERMINAL v8.5.2 "
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(titleView, " Ypsilon 14 | Company Property: Unauthorized access, modification, or use is strictly prohibited.")
	}
	if blinkView, err := g.SetView("blink", maxX-6, 1, maxX-2, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(blinkView, "")
	}
	if menuView, err := g.SetView("menu", 2, 4, maxX/3, (maxY/5)*4-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		menuView.Title = " MENU "
		menuView.Highlight = true
		menuView.SelBgColor = gocui.ColorYellow
		menuView.SelFgColor = gocui.ColorBlack

		for _, menuItem := range MenuItems {
			fmt.Fprintf(menuView, " > %s\n", menuItem)
		}

		if _, err := g.SetCurrentView("menu"); err != nil {
			return err
		}

	}
	if keyView, err := g.SetView("key", 2, (maxY/5)*4-2, maxX/3, (maxY/5)*4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		keyView.Title = " NAVIGATION "
		fmt.Fprintln(keyView, " ENTER=Select    BACKSPACE=Exit")

		if _, err := g.SetCurrentView("menu"); err != nil {
			return err
		}

	}
	if logView, err := g.SetView("log", 2, (maxY/5)*4+1, maxX-2, maxY-1); err != nil {
		logView.Title = " SYSTEM MESSAGES "
		logView.Autoscroll = true
		logView.Wrap = true
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if mainView, err := g.SetView("main", maxX/3+2, 4, maxX-2, (maxY/5)*4); err != nil {
		mainView.Highlight = false
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	// SubViews that show up in certain cases
	if SubViews.Get("legendView") {
		if legendView, err := g.SetView("legendView", (maxX/5)*4+2, 4, maxX-2, (maxY/5)*4); err != nil {
			var legend = []string{
				"ROSTER",
				" 1  Sonya",
				" 2  Ashraf",
				" 3  Dana",
				" 4  Jerome",
				" 5  Kantaro",
				" 6  Morgan",
				" 7  Rie",
				" 8  Rose",
				" 9  Mike",
				" 0  N/A",
				" ",
				"LEGEND",
				" X  Airlock",
				" /  Door",
				"[ ] Elevator",
				" ±  Terminal",
			}

			for _, line := range legend {
				fmt.Fprintf(legendView, " %s\n", line)
			}
			if err != gocui.ErrUnknownView {
				return err
			}

		}
	}

	if SubViews.Get("operationView") {
		if _, err := g.SetView("operationView", (maxX/3)*2+2, 4, maxX-2, (maxY/5)*4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}

	}
	return nil
}
