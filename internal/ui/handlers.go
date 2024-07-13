package ui

import (
	"fmt"
	"log"
	"time"

	"math/rand"

	"github.com/jroimartin/gocui"
)

func WriteToScreen(g *gocui.Gui, v *gocui.View, content []string) {
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

var Controls = []string{" ", "OFF", "OFF", "OFF", "N/A", "OFF", " ", " ", "CLOSED", "CLOSED", "CLOSED", " ", " ", "OPEN", " ", " ", "DEACTIVATED", " "}
var ControlMap = map[string]string{
	"OFF":         "ON",
	"ON":          "OFF",
	"N/A":         "N/A",
	"CLOSED":      "OPEN",
	"OPEN":        "CLOSED",
	"DEACTIVATED": "INITIATED",
	"INITIATED":   "INITIATED",
	" ":           " ",
}

func ToggleControls(ID int) {
	Controls[ID] = ControlMap[Controls[ID]]
}

func HandleOperationView(g *gocui.Gui) error {
	for {
		g.Update(func(g *gocui.Gui) error {
			operationView, err := g.View("operationView")
			if err != nil {
				log.Fatal("Error in HandleOperationView")
				return err
			}
			operationView.Clear()
			for _, line := range Controls {
				fmt.Fprintf(operationView, " %s\n", line)
			}

			return nil
		})
		time.Sleep(10 * time.Millisecond)
		if !SubViews.Get("operationView") {
			return nil
		}
	}
}
