package ui

import (
	"erik-schuetze/ypsilon-14/pkg/util"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

const (
	titleText = " Ypsilon 14 | Company Property: Unauthorized access, modification, or use is strictly prohibited."
	keyText   = " ENTER=Select    BACKSPACE=Exit"
)

var initialValues = map[string]bool{
	"legendView":    false,
	"operationView": false,
}
var SubViews = util.NewMutexMap(initialValues)

var MenuItems = []string{"Station Overview", "Docking Bay", "Floorplan", "Base Operations", "Scan Admin ID Card"}

// -1 = not rendered yet, 0 = one line, 1 = wrapped over two lines.
var (
	titleRenderedMode  = -1
	titleRenderedWidth = -1
	keyRenderedMode    = -1
)

func InitLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	lay, ok := ComputeLayout(maxX, maxY)
	if !ok {
		enterSizeWarning(g)
		renderSizeWarning(g, maxX, maxY)
		return nil
	}

	restoreFocus := ""
	wasWarning := sizeWarningActive
	if wasWarning {
		restoreFocus = leaveSizeWarning(g)
	}

	if titleView, isNew := setView(g, "title", lay.Title); titleView != nil {
		if isNew {
			titleView.Title = " MINING STATION TERMINAL v8.5.2 "
			titleRenderedMode = -1
		}
		renderTitle(titleView, lay.TitleH == 3)
	}
	if blinkView, isNew := setView(g, "blink", lay.Blink); blinkView != nil && isNew {
		fmt.Fprintln(blinkView, "")
	}
	if menuView, isNew := setView(g, "menu", lay.Menu); menuView != nil && isNew {
		menuView.Title = " MENU "
		menuView.Highlight = true
		menuView.SelBgColor = gocui.ColorYellow
		menuView.SelFgColor = gocui.ColorBlack

		for _, menuItem := range MenuItems {
			fmt.Fprintf(menuView, " > %s\n", menuItem)
		}
	}
	if keyView, isNew := setView(g, "key", lay.Key); keyView != nil {
		if isNew {
			keyView.Title = " NAVIGATION "
			keyRenderedMode = -1
		}
		renderKey(keyView, lay.KeyH == 3)
	}
	if logView, isNew := setView(g, "log", lay.Log); logView != nil && isNew {
		logView.Title = " SYSTEM MESSAGES "
		logView.Autoscroll = true
		logView.Wrap = true
	}
	if mainView, isNew := setView(g, "main", lay.Main); mainView != nil && isNew {
		mainView.Highlight = false
	}

	// SubViews that show up in certain cases
	if SubViews.Get("legendView") {
		if legendView, isNew := setView(g, "legendView", lay.Legend); legendView != nil && isNew {
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
		}
	}

	if SubViews.Get("operationView") {
		_, _ = setView(g, "operationView", lay.Operation)
	}

	recenterMessage(g)

	if wasWarning {
		target := restoreFocus
		if target == "" || target == sizeWarningViewName {
			target = "menu"
		}
		if _, err := g.SetCurrentView(target); err != nil {
			_, _ = g.SetCurrentView("menu")
		}
	} else if g.CurrentView() == nil {
		_, _ = g.SetCurrentView("menu")
	}
	return nil
}

// setView updates or creates a view and reports whether it was newly created.
// ComputeLayout guarantees valid rectangles; anything other than ErrUnknownView
// is logged rather than silently leaving a view missing.
func setView(g *gocui.Gui, name string, r Rect) (*gocui.View, bool) {
	v, err := g.SetView(name, r.X0, r.Y0, r.X1, r.Y1)
	switch err {
	case nil:
		return v, false
	case gocui.ErrUnknownView:
		return v, true
	default:
		log.Printf("ui: SetView(%q): %v", name, err)
		return nil, false
	}
}

func renderTitle(v *gocui.View, wrapped bool) {
	mode := 0
	width := -1
	if wrapped {
		mode = 1
		width, _ = v.Size()
	}
	if titleRenderedMode == mode && titleRenderedWidth == width {
		return
	}
	titleRenderedMode = mode
	titleRenderedWidth = width

	v.Clear()
	if !wrapped {
		fmt.Fprintln(v, titleText)
		return
	}

	for _, line := range wrapWords(titleText, width) {
		fmt.Fprintln(v, line)
	}
}

func renderKey(v *gocui.View, wrapped bool) {
	mode := 0
	if wrapped {
		mode = 1
	}
	if keyRenderedMode == mode {
		return
	}
	keyRenderedMode = mode

	v.Clear()
	if !wrapped {
		fmt.Fprintln(v, keyText)
		return
	}
	fmt.Fprintln(v, " ENTER=Select")
	fmt.Fprintln(v, " BACKSPACE=Exit")
}
