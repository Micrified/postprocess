package main

import (
	"os"
	"fmt"
	"io"
	"bufio"
)

type Event struct {
	Executor  int
	Chain     int
	Callback  int
	Start     int
	Duration  int
}

func parse_event (line []byte) *Event {
	split := 0
	var e Event

	// Drop characters until you hit an opening brace
	for _, b := range line {
		if b == '{' {
			break
		}
		split++
	}

	// Check not exceeded line length
	if split >= (len(line) - 1) {
		return nil
	}

	// Construct the format string
	s := string(line[split:])


	// Attempt to parse the arguments
	n, err := fmt.Sscanf(s, "{executor: %d, chain: %d, callback: %d, start: %d, duration: %d}",
		&(e.Executor), &(e.Chain), &(e.Callback), &(e.Start), &(e.Duration))
	if nil != err {
		panic(err)
	}

	// Check if successful
	if n != 5 {
		return nil
	}

	return &e

}

func detect_number_of_chains (events []*Event{})

func main () {
	var filename string = ""
	var file *os.File = nil
	var err error = nil
	var events []*Event = []*Event{}

	// Check arguments
	if len(os.Args) != 2 {
		panic("Required format: " + os.Args[0] + " <filename>")
	} else {
		filename = os.Args[1]
	}

	// Attempt to open the file
	if file, err = os.Open(filename); err != nil {
		panic("Unable to open: " + filename + ": " + err.Error())
	} else {
		defer file.Close()
	}

	// Create a buffered reader
	reader := bufio.NewReader(file)

	for line_num := 1; ; line_num++ {
		line, prefix, err := reader.ReadLine()
		if nil != err {
			break
		}
		if prefix {
			panic("Unable to hand lines that exceed the buffer!")
		}
		//fmt.Printf("%d. %s", line_num, line)
		e := parse_event(line)
		if e != nil {
			events = append(events, e)
		} else {
			fmt.Printf("No event for line: " + string(line) + "\n")
		}
	}

	if err != nil && err != io.EOF {
		panic("Error while reading lines: " + err.Error())
	}

	// Return early if nothing to do
	if len(events) == 0 {
		fmt.Printf("No events ...\n")
		return
	}

	// Adjust the time by subtracting the minimum time (likely the first)
	min_time := events[0].Start
	for _, e := range events {
		if e.Start < min_time {
			min_time = e.Start
		}
	}
	for _, e := range events {
		e.Start -= min_time
	}

	// Output the results
	for _, e := range events {
		fmt.Printf("%d, %d, %d, %d, %d\n", e.Executor, e.Chain, e.Callback, e.Start, e.Duration)
	}
}