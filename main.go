package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jroimartin/gocui"
)

var LogRoutineRunning = false
var BlinkRoutineRunning = false
var ControlRoutineRunning = false
var ControlViewActive = false
var LockControls = false
var ViewsToRemove = []string{}
var MenuItems = []string{"Station Overview", "Docking Bay History", "Floorplan", "Base Operations"}
var Controls = []string{" ", "OFF", "OFF", "OFF", "OFF", "OFF", " ", " ", "CLOSED", "CLOSED", "CLOSED", " ", " ", "DEACTIVATED"}
var ControlMap = map[string]string{
	"OFF":         "ON",
	"ON":          "OFF",
	"CLOSED":      "OPEN",
	"OPEN":        "CLOSED",
	"DEACTIVATED": "INITIATED",
	"INITIATED":   "INITIATED",
	" ":           " ",
}

/*
TODO

PIN for self destruction / ID CARD

think about additional Info / Flavor / Controls / Menu Options
*/

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := initKeyBindings(g); err != nil {
		log.Panicln(err)
	}

	// start one concurrent thread that handles writing Log lines to the System Messages
	if !LogRoutineRunning {
		LogRoutineRunning = true
		go writeLogLines(g)
	}

	// start one concurrent thread for the small Logo Box in the upper right corner
	if !BlinkRoutineRunning {
		LogRoutineRunning = true
		go blinkEffect(g)
	}

	// start one concurrent thread for the small Logo Box in the upper right corner
	if !ControlRoutineRunning {
		ControlRoutineRunning = true
		go handleControlView(g)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func blinkEffect(g *gocui.Gui) error {
	blinkString := ""
	for {
		g.Update(func(g *gocui.Gui) error {
			blinkView, err := g.View("blink")
			if err != nil {
				return err
			}
			blinkView.Clear()
			fmt.Fprintf(blinkView, "%s", blinkString)

			switch blinkString {
			case "":
				blinkString = "Y  "
			case "Y  ":
				blinkString = "Y1 "
			case "Y1 ":
				blinkString = "Y14"
			case "Y14":
				blinkString = ""

			}

			return nil
		})

		randomTimeInterval := rand.Intn(600-300) + 300
		time.Sleep(time.Duration(randomTimeInterval) * time.Millisecond)
	}

	return nil
}

func writeLogLines(g *gocui.Gui) error {
	logMessages := []string{
		"INFO:    All core systems operating within normal parameters.",
		"INFO:    No external weather anomalies detected. Station shields at standard levels.",
		"INFO:    Supply Shipment arrived on Docking Bay 2.",
		"INFO:    Routine health check complete. No critical issues found.",
		"INFO:    Excavator drones operating at 95.3% efficiency.",
		"WARNING: Detected pressure anomaly in water supply.",
		"INFO:    50t of ore extracted, exceeding the weekly target by 2.7%.",
		"INFO:    Ore Storage Capacity at 43.3%. Cargo transfer scheduled in 15 days.",
		"INFO:    All airlocks secure. No breach detected in the last 24 hours.",
		"INFO:    Routine security sweep initiated. Access logs being reviewed.",
		"WARNING: Unstable power supply. Power surge detected in elevator shaft.",
		"INFO:    19/20 mining drills deployed.",
		"WARNING: Airlock 1 Override active",
		"INFO:    Routine maintenance of mining equipment completed. All systems operational.",
		"INFO:    Atmospheric pressure levels in all sectors normalized within standard tolerances.",
		"INFO:    Primary power systems operating stable.",
		"INFO:    Gravity generators operating at 1.0G. No deviations reported in the last 24 hours.",
		"INFO:    Temperature and Humidity levels remain within standard limits.",
		"INFO:    Running check on backup systems for routine diagnostics.",
	}

	for {
		for _, entry := range logMessages {
			g.Update(func(g *gocui.Gui) error {
				logView, err := g.View("log")
				if err != nil {
					return err
				}
				fmt.Fprintf(logView, "\n %s", entry)
				return nil
			})
			randomTimeInterval := rand.Intn(1400-400) + 400
			time.Sleep(time.Duration(randomTimeInterval) * time.Millisecond)
		}
	}
	return nil
}

func layout(g *gocui.Gui) error {
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
	return nil
}

func initKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("menu", gocui.KeyEnter, gocui.ModNone, selectMenuItem); err != nil {
		return err
	}
	if err := g.SetKeybinding("menu", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("menu", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, selectMenuItem); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyBackspace2, gocui.ModNone, exitMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func selectMenuItem(g *gocui.Gui, activeView *gocui.View) error {
	if !LockControls {
		menuView, err := g.View("menu")
		if err != nil {
			return err
		}
		_, menuY := menuView.Cursor()

		if activeView == nil || activeView.Name() == "menu" {
			mainView, err := g.View("main")
			if err != nil {
				return err
			}

			switch menuY {
			case 0:
				stationOverview(g, mainView)
			case 1:
				dockingBayHistory(g, mainView)
			case 2:
				floorplan(g, mainView)
			case 3:
				baseOperations(g, mainView)
			default:
				return nil
			}

			g.Cursor = false
			_, err = g.SetCurrentView("main")
			mainView.SetOrigin(0, 0)
			mainView.SetCursor(mainView.Origin())
			return err
		} else if activeView.Name() == "main" && menuY == 3 {
			_, mainY := activeView.Cursor()
			if mainY < 14 {
				toggleControls(mainY)
			}
		}
	}
	return nil
}

func toggleControls(ID int) {
	Controls[ID] = ControlMap[Controls[ID]]
}

func writeToScreen(g *gocui.Gui, v *gocui.View, content []string) {
	LockControls = true
	for _, line := range content {
		fmt.Fprintf(v, " ")
		for _, char := range line {
			g.Update(func(g *gocui.Gui) error {
				fmt.Fprintf(v, "%c", char)
				return nil
			})
			randomTimeInterval := rand.Intn(10-1) + 1
			time.Sleep(time.Duration(randomTimeInterval) * time.Millisecond)
		}
		fmt.Fprintf(v, "\n")
	}
	LockControls = false
}

func stationOverview(g *gocui.Gui, v *gocui.View) error {
	overview := []string{
		"STATION NAME:     Ypsilon 14",
		"STATION TYPE:     Asteroid Mining Station",
		"MINING MODULE:    OrionMiningCorp X7 rev 2.1",
		"QUARTERS MODULE:  OrionMiningCorp E2 rev 1.1",
		"OVERALL STATUS:   normal",
	}
	go writeToScreen(g, v, overview)
	return nil
}

func dockingBayHistory(g *gocui.Gui, v *gocui.View) error {
	history := []string{
		"DOCKING BAY LOG:",
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
	}
	v.Autoscroll = true
	go writeToScreen(g, v, history)
	return nil
}

func floorplan(g *gocui.Gui, v *gocui.View) error {
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
		" |     ===     |  | 5 /              /      |",
		" +-------------+  +---+   +-/-+-/-+--+ Wash |",
		" 	     |v|        | 4 /   | 3 | 2 |  | room |",
		" 	     |v|        +---+---+---+---+  +------+",
		" +-------------+",
		" |  Mineshaft  X . . . ",
		" +-------------+",
	}
	legend := []string{
		"ROSTER",
		"1 - Sonya",
		"2 - Ashraf",
		"3 - Dana",
		"4 - Jerome",
		"5 - Kantaro",
		"6 - Morgan",
		"7 - Rie",
		"8 - Rose",
		"9 - Mike",
		"0 - N/A",
		" ",
		"LEGEND",
		"X - Airlock",
		"/ - Door",
		"= - Elevator",
		"± - Terminal",
	}

	// print floorplan in main view
	go writeToScreen(g, v, plan)

	// create legend gocui.View to show them on the right
	maxX, maxY := g.Size()
	legendView, err := g.SetView("legendView", (maxX/5)*4+2, 4, maxX-2, (maxY/5)*4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	for _, line := range legend {
		fmt.Fprintf(legendView, " %s\n", line)
	}
	//go writeToScreen(g, legendView, legend)

	// add crewView to list of views to remove when exiting the mainView
	ViewsToRemove = append(ViewsToRemove, "legendView")

	return nil
}

func baseOperations(g *gocui.Gui, v *gocui.View) error {
	operations := []string{
		"SHOWER CONTROL",
		"> Shower 1",
		"> Shower 2",
		"> Shower 3",
		"> Shower 4",
		"> Shower 5",
		" ",
		"AIRLOCK CONTROL",
		"> Airlock 1",
		"> Airlock 2",
		"> Airlock 3",
		" ",
		"CRITICAL OPERATIONS",
		"> Initiate Self Destruction",
	}

	v.Highlight = true
	v.SelBgColor = gocui.ColorYellow
	v.SelFgColor = gocui.ColorBlack

	go writeToScreen(g, v, operations)

	ControlViewActive = true

	// create legend gocui.View to show them on the right
	maxX, maxY := g.Size()
	_, err := g.SetView("control", (maxX/3)*2+2, 4, maxX-2, (maxY/5)*4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	ViewsToRemove = append(ViewsToRemove, "control")

	return nil
}

func handleControlView(g *gocui.Gui) error {
	for {
		if ControlViewActive {
			g.Update(func(g *gocui.Gui) error {
				controlView, err := g.View("control")
				if err != nil {
					return err
				}
				controlView.Clear()
				for _, line := range Controls {
					fmt.Fprintf(controlView, " %s\n", line)
				}

				return nil
			})
			time.Sleep(5 * time.Millisecond)
		}
	}
	return nil
}

func exitMain(g *gocui.Gui, v *gocui.View) error {
	if !LockControls {
		if ControlViewActive {
			ControlViewActive = false
		}
		if v.Name() == "main" {
			v.Clear()
			g.Cursor = true
			v.Highlight = false
			v.Autoscroll = false
			for _, view := range ViewsToRemove {
				g.DeleteView(view)
			}
			_, err := g.SetCurrentView("menu")
			return err
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}
