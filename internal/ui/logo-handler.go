package ui

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jroimartin/gocui"
)

func StartLogoAnimation(g *gocui.Gui) error {
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
}
