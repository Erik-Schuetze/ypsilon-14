package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jroimartin/gocui"
)

var LogRoutineRunning = false

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keyBindings(g); err != nil {
		log.Panicln(err)
	}

	if !LogRoutineRunning {
		LogRoutineRunning = true
		go writeLogLines(g)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
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
			randomTimeInterval := rand.Intn(1200-400) + 400
			time.Sleep(time.Duration(randomTimeInterval) * time.Millisecond)
		}
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if titleView, err := g.SetView("title", 2, 1, maxX-2, 3); err != nil {
		titleView.Title = " MINING STATION TERMINAL v8.5.2 "
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(titleView, " Name: Ypsilon 14 | Type: Asteroid Mining Station | Status: Healthy | Team Leader: Sonya")
	}
	if menuView, err := g.SetView("menu", 2, 4, maxX/3, (maxY/5)*4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		menuView.Title = " MENU "
		menuView.Highlight = true
		menuView.SelBgColor = gocui.ColorYellow
		menuView.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(menuView, "Station Overview")
		fmt.Fprintln(menuView, "Cargo History")
		fmt.Fprintln(menuView, "Base Operations")
		fmt.Fprintln(menuView, "Critical Operations (locked)")

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
		mainView.Highlight = true
		mainView.Highlight = true
		mainView.SelBgColor = gocui.ColorYellow
		mainView.SelFgColor = gocui.ColorBlack
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

func keyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("menu", gocui.KeyEnter, gocui.ModNone, selectMenuItem); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyBackspace2, gocui.ModNone, exitMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("menu", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("menu", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func selectMenuItem(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "menu" {
		_, err := g.SetCurrentView("main")
		return err
	}
	return nil
}

func exitMain(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == "main" {
		_, err := g.SetCurrentView("menu")
		return err
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
