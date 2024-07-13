package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

var LockControls = false

func InitKeyBindings(g *gocui.Gui) error {
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
	if err := g.SetKeybinding("message", gocui.KeyBackspace2, gocui.ModNone, exitMain); err != nil {
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
			case 0: // Station Overview
				InitStationOverview(g, mainView)
			case 1: // Docking Bay History
				InitDockingBayHistory(g, mainView)
			case 2: // Floorplan
				InitFloorplan(g, mainView)
			case 3: // Base Operations
				InitBaseOperations(g, mainView)
			case 4: // Scan ID Card
				AdminAccess = true
				DisplayMessage(g, "CARD READER v1.2", "Enabled admin access")
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
			if mainY < 16 {
				ToggleControls(mainY)
			} else if mainY == 16 {
				if AdminAccess {
					ToggleControls(mainY)
					DisplayMessage(g, "SELF DESTRUCTION", "Initiated Self Destruction. Evacuate immediately!")
				} else {
					DisplayMessage(g, "ERROR", "Admin access required!")
				}
			}
		}
	}
	return nil
}

func exitMain(g *gocui.Gui, v *gocui.View) error {
	if !LockControls {
		if ControlViewActive {
			ControlViewActive = false
		}
		if v.Name() == "main" || v.Name() == "message" {
			v, _ := g.View("main")
			v.Clear()
			g.Cursor = true
			v.Highlight = false
			v.Autoscroll = false
			v.Editable = false
			for _, view := range ViewsToRemove {
				g.DeleteView(view)
				SubViews.Write(view, false)
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
