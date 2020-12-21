package main

import (

	// Standard packages
	"fmt"
	"os"
	"sort"
	
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

	// Sort the log file chronologically
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Start_us < events[j].Start_us
	})

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

	// Normalize all results, and append it to the given file
	// Results are 1:1 with chains
	// chain_id chain_count rand_seed mode bcrt wcrt acrt util period prio
	for i, result := range results {

		b2s := func (x bool) int {
			if x {
				return 1
			}
			return 0
		}

		chain_id    := result.ID
		chain_count := len(chains)
		chain_len   := len(chains[i].Path)
		rand_seed   := chains[i].Random_seed
		mode        := b2s(chains[i].PPE)
		period      := chains[i].Period_us
		util        := chains[i].Utilisation
		bcrt        := result.BCRT_us
		wcrt        := result.WCRT_us
		acrt        := result.ACRT_us
		fmt.Printf("%d %d %d %d %d %d %f %d %d %d\n",
			chain_id, chain_count, chain_len, rand_seed, mode, period, util, bcrt, wcrt, acrt)
	}

}