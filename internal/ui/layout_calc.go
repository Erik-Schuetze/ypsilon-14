package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	MinTerminalWidth  = 100
	MinTerminalHeight = 29

	MenuTextWidth       = 21
	LegendTextWidth     = 13
	OperationTextWidth  = 12
	MainTextWidth       = 49 // station overview incl. leading space
	FloorplanWidth      = 46 // rendered incl. leading space
	FloorplanHeight     = 18
	LegendHeight        = 17
	OperationHeight     = 18
	OperationsListWidth = 33 // longest operations line incl. leading space
)

// Derived from the user-visible strings in layout.go so they cannot desync.
var (
	TitleTextWidth = runeLen(titleText) // includes leading space
	KeyTextWidth   = runeLen(keyText)
)

// Rect is a gocui view rectangle. InnerW/InnerH match gocui's View.Size().
type Rect struct {
	X0, Y0, X1, Y1 int
}

func (r Rect) InnerW() int { return r.X1 - r.X0 - 1 }
func (r Rect) InnerH() int { return r.Y1 - r.Y0 - 1 }

// Valid reports whether gocui can create a usable view for this rectangle.
func (r Rect) Valid(width, height int) bool {
	return r.X0 >= 0 && r.Y0 >= 0 && r.X1 < width && r.Y1 < height &&
		r.InnerW() > 0 && r.InnerH() > 0
}

type Layout struct {
	Title, Blink, Menu, Key, Log, Main, Legend, Operation Rect
	ContentTop, TitleH, KeyH                              int
}

// ComputeLayout returns the view geometry for a terminal size. The second
// return value is false when the terminal is too small for the fixed content.
func ComputeLayout(width, height int) (Layout, bool) {
	if width < MinTerminalWidth || height < MinTerminalHeight {
		return Layout{}, false
	}

	// The title body is one line when it fits, otherwise two word-wrapped lines.
	titleH := 2
	if width-10 < TitleTextWidth {
		titleH = 3
	}
	contentTop := titleH + 2

	// splitY keeps 80% of the rows for the main area, retaining the remainder
	// in the main area instead of silently dropping it (monotonic in height).
	splitY := height - height/5

	// The navigation hint needs a second line when the menu column is narrow.
	keyH := 2
	if width/3-3 < KeyTextWidth {
		keyH = 3
	}

	mainX0 := width/3 + 2

	// Split the space right of the main view between the visible main content
	// and an overlay. Both sides then grow monotonically with the terminal.
	overlayX0 := func(minVisible int) int {
		visible := (width - mainX0) / 2
		if visible < minVisible {
			visible = minVisible
		}
		return mainX0 + 1 + visible
	}

	l := Layout{
		Title:      Rect{2, 1, width - 7, 1 + titleH},
		Blink:      Rect{width - 6, 1, width - 2, 3},
		Menu:       Rect{2, contentTop, width / 3, splitY - keyH - 1},
		Key:        Rect{2, splitY - keyH, width / 3, splitY},
		Log:        Rect{2, splitY + 1, width - 2, height - 1},
		Main:       Rect{mainX0, contentTop, width - 2, splitY},
		Legend:     Rect{overlayX0(FloorplanWidth), contentTop, width - 2, splitY},
		Operation:  Rect{overlayX0(OperationsListWidth), contentTop, width - 2, splitY},
		ContentTop: contentTop,
		TitleH:     titleH,
		KeyH:       keyH,
	}

	// Never hand gocui a rectangle it would reject.
	for _, r := range []Rect{l.Title, l.Blink, l.Menu, l.Key, l.Log, l.Main, l.Legend, l.Operation} {
		if !r.Valid(width, height) {
			return Layout{}, false
		}
	}
	return l, true
}

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

// wrapWords greedily wraps text on spaces. Every line keeps one leading space,
// matching the original content style.
func wrapWords(text string, width int) []string {
	if width < 1 {
		return []string{text}
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0, len(words))
	cur := " " + words[0]
	for _, word := range words[1:] {
		candidate := cur + " " + word
		if runeLen(candidate) > width {
			lines = append(lines, cur)
			cur = " " + word
		} else {
			cur = candidate
		}
	}
	return append(lines, cur)
}

// wrappedLineCount mirrors gocui's cell-based wrapping for the ASCII content
// used by this application.
func wrappedLineCount(text string, innerW int) int {
	if innerW < 1 {
		return 1
	}
	total := 0
	for _, line := range strings.Split(text, "\n") {
		n := runeLen(line)
		if n < innerW {
			total++
		} else {
			// gocui's loop (n <= len) emits a trailing empty line when the
			// line length is an exact multiple of the inner width.
			total += n/innerW + 1
		}
	}
	if total == 0 {
		return 1
	}
	return total
}

// messageRect computes a centered, screen-clamped popup rectangle and the
// inner width it provides.
func messageRect(screenW, screenH int, title, message string) (Rect, int) {
	if screenW < 4 || screenH < 4 {
		return Rect{}, 0
	}

	body := " " + message
	maxInnerW := screenW - 6
	if maxInnerW < 1 {
		maxInnerW = 1
	}

	// gocui's drawTitle only renders innerW-2 runes, and the title is
	// rendered as " "+title+" ", so the title needs N+4 inner columns.
	wantW := runeLen(body)
	if titleW := runeLen(" "+title+" ") + 2; titleW > wantW {
		wantW = titleW
	}
	innerW := wantW
	if innerW > maxInnerW {
		innerW = maxInnerW
	}

	innerH := wrappedLineCount(body, innerW)
	maxInnerH := screenH - 4
	if maxInnerH < 1 {
		maxInnerH = 1
	}
	if innerH > maxInnerH {
		innerH = maxInnerH
	}

	x0 := (screenW - (innerW + 1)) / 2
	y0 := (screenH - (innerH + 1)) / 2
	return Rect{x0, y0, x0 + innerW + 1, y0 + innerH + 1}, innerW
}

// sizeWarningText returns centered text for the too-small-terminal screen.
func sizeWarningText(w, h int) string {
	if w < 3 || h < 3 {
		return ""
	}

	lines := []string{
		"TERMINAL TOO SMALL",
		fmt.Sprintf("Required: %d x %d", MinTerminalWidth, MinTerminalHeight),
		fmt.Sprintf("Current:  %d x %d", w, h),
		"",
		"Enlarge the terminal window.",
		"CTRL-C to quit.",
	}

	innerW := w - 2
	innerH := h - 2
	maxLen := 0
	for _, line := range lines {
		if n := runeLen(line); n > maxLen {
			maxLen = n
		}
	}

	left := 0
	if innerW > maxLen {
		left = (innerW - maxLen) / 2
	}
	top := 0
	if innerH > len(lines) {
		top = (innerH - len(lines)) / 2
	}

	var b strings.Builder
	for i := 0; i < top; i++ {
		b.WriteByte('\n')
	}
	for _, line := range lines {
		if left > 0 {
			b.WriteString(strings.Repeat(" ", left))
		}
		if innerW > 0 && runeLen(line) > innerW {
			line = string([]rune(line)[:innerW])
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}
