package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const sizeWarningViewName = "sizeWarning"

var (
	sizeWarningActive  bool
	focusBeforeWarning string
	warnedW, warnedH   int
)

// enterSizeWarning records the focused view so it can be restored after the
// terminal grows back to a supported size.
func enterSizeWarning(g *gocui.Gui) {
	if sizeWarningActive {
		return
	}
	focusBeforeWarning = ""
	if cv := g.CurrentView(); cv != nil {
		focusBeforeWarning = cv.Name()
	}
	sizeWarningActive = true
}

func renderSizeWarning(g *gocui.Gui, w, h int) {
	// A 1- or 2-cell terminal cannot hold a framed view. Do nothing rather
	// than calling SetView with a degenerate rectangle.
	if w < 3 || h < 3 {
		return
	}

	v, err := g.SetView(sizeWarningViewName, 0, 0, w-1, h-1)
	if err != nil && err != gocui.ErrUnknownView {
		return
	}
	if err == gocui.ErrUnknownView {
		v.Title = " TERMINAL "
		v.Highlight = false
		// Wrap must remain false: gocui errors when drawing a wrapping view
		// whose inner width is zero.
	}

	_, _ = g.SetViewOnTop(sizeWarningViewName)

	if err == gocui.ErrUnknownView || warnedW != w || warnedH != h {
		warnedW, warnedH = w, h
		v.Clear()
		fmt.Fprint(v, sizeWarningText(w, h))
	}

	if cv := g.CurrentView(); cv == nil || cv.Name() != sizeWarningViewName {
		_, _ = g.SetCurrentView(sizeWarningViewName)
	}
}

// leaveSizeWarning removes the warning and returns the view that had focus
// before the terminal became too small.
func leaveSizeWarning(g *gocui.Gui) string {
	_ = g.DeleteView(sizeWarningViewName)
	sizeWarningActive = false
	warnedW, warnedH = 0, 0
	target := focusBeforeWarning
	focusBeforeWarning = ""
	return target
}
