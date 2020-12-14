package main

import (

	// Standard packages
	"fmt"
	"os"

	// Custom packages
	"analysis"

	// Third party packages
	"github.com/gookit/color"
)


func warn (s string, args ...interface{}) {
	color.Style{color.FgYellow, color.OpBold}.Printf("%s\n", fmt.Sprintf(s, args...))
}

func info (s string, args ...interface{}) {
	color.Style{color.FgGreen, color.OpBold}.Printf("%s\n", fmt.Sprintf(s, args...))
}

func main () {
	usage_fmt := "%s <chains.json> <logfile>"

	// Verify program arguments
	if len(os.Args) != 3 {
		warn(usage_fmt, os.Args[0])
		return
	}

	// Attempt to open the chains file
	chains, err := analysis.ReadChains(os.Args[1])
	if nil != err {
		panic(err.Error())
	}

	// Attempt to open the log file, and convert it to events
	events, err := analysis.ReadEvents(os.Args[2])
	if nil != err {
		panic(err.Error())
	}

	// Output the chain information
	warn("---- Expected chains (from JSON file) ----")
	for _, chain := range chains {
		fmt.Printf("Chain %d: Prio: %d, Path: %s, Period: %dus, Utilisation: %f\n",
			chain.ID, chain.Prio, analysis.Path2String(chain.Path), chain.Period_us,
			chain.Utilisation)
	}

	// Compute analysis
	warn("---- Analyzing %d events ----", len(events))
	results := analysis.Analyze(chains, events)
	warn("---- Results ----")
	for _, result := range results {
		fmt.Printf("Chain %d: WCRT: %dus, ACRT: %dus, BCRT: %dus\n", result.ID, 
			result.WCRT_us, result.ACRT_us, result.BCRT_us)
	}

}