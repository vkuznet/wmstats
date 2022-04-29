package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

// git version of our code
var version string

func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("wmstats git=%s go=%s date=%s", version, goVersion, tstamp)
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var input string
	flag.StringVar(&input, "input", "", "input file")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbosity level")
	flag.Parse()
	if version {
		fmt.Println("wmstats version:", info())
		return
	}
	server(input, verbose)
}
