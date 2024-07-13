package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var ViewsToRemove = []string{}
var ControlViewActive = false
var AdminAccess = false

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
	ViewsToRemove = append(ViewsToRemove, "legendView")

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

	ViewsToRemove = append(ViewsToRemove, "operationView")

	return nil
}

func DisplayMessage(g *gocui.Gui, title string, message string) error {
	maxX, maxY := g.Size()
	width := len(message) + 2
	messageView, err := g.SetView("message", maxX/2-width/2, maxY/2-1, maxX/2+width/2+1, maxY/2+1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	g.SetViewOnTop("message")
	g.SetCurrentView("message")
	messageView.Title = fmt.Sprintf(" %s ", title)
	fmt.Fprintf(messageView, " %s", message)

	ViewsToRemove = append(ViewsToRemove, "message")

	return nil
}
