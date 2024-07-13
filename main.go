package main

import (
	"erik-schuetze/ypsilon-14/internal/ui"
	"log"

	"github.com/jroimartin/gocui"
)

var LogRoutineRunning = false
var BlinkRoutineRunning = false
var ControlRoutineRunning = false

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(ui.InitLayout)

	err = ui.InitKeyBindings(g)
	if err != nil {
		log.Panicln(err)
	}

	// start one concurrent thread that handles writing Log lines to the System Messages
	if !LogRoutineRunning {
		LogRoutineRunning = true
		go ui.StartLoglineAnimation(g)
	}

	// start one concurrent thread for the small Logo Box in the upper right corner
	if !BlinkRoutineRunning {
		BlinkRoutineRunning = true
		go ui.StartLogoAnimation(g)
	}

	// start one concurrent thread for handling the controls in the base operation menu
	if !ControlRoutineRunning {
		ControlRoutineRunning = true
		//go ui.HandleOperationView(g)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
