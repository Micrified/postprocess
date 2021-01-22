package main

import (

	// Standard packages
	"fmt"
	"os"
	"sort"
	
	// Custom packages
	"analysis"
	"types"

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
		for i := 1; i < len(events); i++ {
			events[i].Start_us -= events[0].Start_us
		}
		events[0].Start_us = 0
	}

	// Output the chain information
	fmt.Fprintf(os.Stderr, "---- Expected chains (from JSON file) ----\n")
	for _, chain := range chains {
		fmt.Fprintf(os.Stderr, "Chain %d: Prio: %d, Path: %s, Period: %dus, Utilisation: %f\n",
			chain.ID, chain.Prio, analysis.Path2String(chain.Path), chain.Period_us,
			chain.Utilisation)
	}

	// Compute analysis
	fmt.Fprintf(os.Stderr, "---- Analysing %d events ----\n", len(events))
	results := analysis.Analyse(chains, events)
	fmt.Fprintf(os.Stderr, "---- Results ----\n")
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

		// Build a trace
		t := types.Trace{
			ID: result.ID,
			Priority: chains[i].Prio,
			Length: len(chains[i].Path),
			Period: chains[i].Period_us,
			Utilisation: chains[i].Utilisation,
			BCRT_us: result.BCRT_us,
			WCRT_us: result.WCRT_us,
			ACRT_us: result.ACRT_us,
			Chain_count: len(chains),
			Avg_chain_length: chains[i].Avg_len,
			Seed: chains[i].Random_seed,
			Merge_p: chains[i].Merge_p,
			Sync_p: chains[i].Sync_p,
			Variance: chains[i].Variance,
			PPE: b2s(chains[i].PPE),
			Executors: chains[i].Executors,
		}

		// Serialise trace to STDOUT
		t.Serialise(os.Stdout)
	}
}