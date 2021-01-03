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

	// To avoid possible overflow errors later, subtract earliest time from all logs
	if len(events) > 0 {
		first := events[0]
		for i := 1; i < len(events); i++ {
			events[i].Start_us -= first.Start_us
		}
		events[0].Start_us = 0
	}

	// Output the chain information
	fmt.Fprintf(os.Stderr, "---- Expected chains (from JSON file) ----")
	for _, chain := range chains {
		fmt.Fprintf(os.Stderr, "Chain %d: Prio: %d, Path: %s, Period: %dus, Utilisation: %f\n",
			chain.ID, chain.Prio, analysis.Path2String(chain.Path), chain.Period_us,
			chain.Utilisation)
	}

	// Compute analysis
	fmt.Fprintf(os.Stderr, "---- Analyzing %d events ----", len(events))
	results := analysis.Analyze(chains, events)
	fmt.Fprintf(os.Stderr, "---- Results ----")
	for _, result := range results {
		fmt.Fprintf(os.Stderr, "Chain %d: WCRT: %dus, ACRT: %dus, BCRT: %dus\n", result.ID, 
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

		// Chain specific attributes
		chain_id      := result.ID
		chain_prio    := chains[i].Prio
		chain_len     := len(chains[i].Path)
		period        := chains[i].Period_us
		util          := chains[i].Utilisation
		bcrt          := result.BCRT_us
		wcrt          := result.WCRT_us
		acrt          := result.ACRT_us

		fmt.Fprintf(os.Stdout, "%d %d %d %d %f %d %d %d ",
			chain_id, chain_prio, chain_len, period, util, bcrt, wcrt, acrt)

		// General test attributes
		chain_count   := len(chains)
		chain_avg_len := chains[i].Avg_len
		rand_seed     := chains[i].Random_seed
		merge_p       := chains[i].Merge_p
		sync_p        := chains[i].Sync_p
		variance      := chains[i].Variance
		mode          := b2s(chains[i].PPE)

		fmt.Fprintf(os.Stdout, "%d %d %d %f %f %f %d\n",
			chain_count, chain_avg_len, rand_seed, merge_p, sync_p, variance, mode)
	}
}