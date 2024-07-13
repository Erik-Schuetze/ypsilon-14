package ui

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jroimartin/gocui"
)

var logMessages = []string{
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
	"INFO:    Bleep Bloop",
	"INFO:    Gravity generators operating at 1.0G. No deviations reported in the last 24 hours.",
	"INFO:    Temperature and Humidity levels remain within standard limits.",
	"INFO:    Running check on backup systems for routine diagnostics.",
}

func StartLoglineAnimation(g *gocui.Gui) error {

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
}
