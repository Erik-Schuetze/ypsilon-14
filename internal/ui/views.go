package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var ViewsToRemove = []string{}
var ControlViewActive = false
var AdminAccess = false

var (
	popupTitle   string
	popupMessage string
)

func addViewToRemove(name string) {
	for _, existing := range ViewsToRemove {
		if existing == name {
			return
		}
	}
	ViewsToRemove = append(ViewsToRemove, name)
}

func InitStationOverview(g *gocui.Gui, v *gocui.View) error {
	overview := []string{
		"STATION NAME:       Ypsilon 14",
		"STATION TYPE:       Asteroid Mining Station",
		" ",
		"MINING MODULE:      OrionMiningCorp X7 rev 2.1",
		"QUARTERS MODULE:    OrionMiningCorp E2 rev 1.1",
		" ",
		"OVERALL STATUS:     normal",
		"COMPLIANCE STATUS:  compliant",
		" ",
		"UPDATES:            14 Updates available",
		"                    unvoluntary reboot scheduled",
	}
	go WriteToScreen(g, v, overview)
	return nil
}

func InitDockingBayHistory(g *gocui.Gui, v *gocui.View) error {
	history := []string{
		"DOCKING BAY LOG:",
		"10 months ago  arrival     cargo transport",
		"10 months ago  departure   cargo transport",
		"9 months ago   arrival     research probe",
		"9 months ago   departure   waste products",
		"9 months ago   arrival     passenger shuttle",
		"9 months ago   departure   passenger shuttle",
		"8 months ago   departure   research probe",
		"8 months ago   arrival     cargo transport",
		"8 months ago   departure   cargo transport",
		"8 months ago   arrival     mining equipment",
		"8 months ago   departure   medical support",
		"7 months ago   arrival     cargo transport",
		"6 months ago   departure   waste products",
		"6 months ago   arrival     mining equipment",
		"6 months ago   departure   technician support",
		"4 months ago   departure   cargo transport",
		"4 months ago   arrival     passenger shuttle",
		"3 months ago   departure   passenger shuttle",
		"3 months ago   departure   technician support",
		"2 months ago   arrival     cargo transport",
		"2 months ago   departure   cargo transport",
		"5 weeks ago    arrival     ********",
		"1 hour ago     arrival     supply transport",
		" ",
		"UPCOMING SHIPMENTS:",
		"in 2 weeks     arrival     drill parts",
		" ",
		"CURRENT BAY STATUS:",
		"BAY 1:         ********",
		"BAY 2:         Tempest",
	}
	v.Autoscroll = true
	go WriteToScreen(g, v, history)
	return nil
}

func InitFloorplan(g *gocui.Gui, v *gocui.View) error {
	plan := []string{
		" ",
		" _____     _____",
		" Bay 1     Bay 2",
		" --+--     --+--",
		"   |         |    +---+---+---+---+  +------+",
		"   |         |    | 7 /   | 8 | 9 |  | Mess |",
		" +-X---------X-+  +---+   +-/-+-/-+--+      |",
		" |      ±      |  | 6 /              /      |",
		" |             +--+---+   +---+---+---------+",
		" |  Workspace             | 0 | 1 |          ",
		" |             +--+---+   +-/-+-/-+---------+",
		" |     [ ]     |  | 5 /              /      |",
		" +-------------+  +---+   +-/-+-/-+--+ Wash |",
		"       |v|        | 4 /   | 3 | 2 |  | room |",
		"       |v|        +---+---+---+---+  +------+",
		" +-------------+",
		" |  Mineshaft  X . . . ",
		" +-------------+",
	}

	// print floorplan in main view
	go WriteToScreen(g, v, plan)

	// set this to show the legend on the right
	SubViews.Write("legendView", true)

	// add legend to list of views to remove when exiting the mainView
	addViewToRemove("legendView")

	return nil
}

func InitBaseOperations(g *gocui.Gui, v *gocui.View) error {
	operations := []string{
		"SHOWER OVERRIDE",
		"> Shower 1",
		"> Shower 2",
		"> Shower 3",
		"> Shower 4",
		"> Shower 5",
		" ",
		"AIRLOCK CONTROL",
		"> Airlock 1 (Bay 1)",
		"> Airlock 2 (Bay 2)",
		"> Airlock 3 (Mineshaft)",
		" ",
		"AIR VENT CONTROL",
		"> Central Air Circulation System",
		" ",
		"CRITICAL OPERATIONS",
		"> Initiate Self Destruction",
	}

	v.Highlight = true
	v.SelBgColor = gocui.ColorYellow
	v.SelFgColor = gocui.ColorBlack

	go WriteToScreen(g, v, operations)

	//ControlViewActive = true

	// create legend gocui.View to show them on the right
	SubViews.Write("operationView", true)
	go HandleOperationView(g)

	addViewToRemove("operationView")

	return nil
}

func DisplayMessage(g *gocui.Gui, title string, message string) error {
	popupTitle, popupMessage = title, message
	layoutMessage(g, true)
	return nil
}

// layoutMessage positions the popup. When rewrite is false it only updates
// geometry (used when the terminal is resized while the popup is open).
func layoutMessage(g *gocui.Gui, rewrite bool) {
	w, h := g.Size()
	if w < 4 || h < 4 {
		return
	}

	r, _ := messageRect(w, h, popupTitle, popupMessage)
	messageView, err := g.SetView("message", r.X0, r.Y0, r.X1, r.Y1)
	if err != nil && err != gocui.ErrUnknownView {
		return
	}

	messageView.Wrap = true
	messageView.Title = fmt.Sprintf(" %s ", popupTitle)
	if rewrite || err == gocui.ErrUnknownView {
		messageView.Clear()
		fmt.Fprintf(messageView, " %s", popupMessage)
	}

	_, _ = g.SetViewOnTop("message")
	addViewToRemove("message")
	if rewrite {
		_, _ = g.SetCurrentView("message")
	}
}

// recenterMessage keeps an open popup centered after a resize.
func recenterMessage(g *gocui.Gui) {
	if popupMessage == "" {
		return
	}
	if _, err := g.View("message"); err != nil {
		return
	}
	layoutMessage(g, false)
}
