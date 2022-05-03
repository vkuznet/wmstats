package main

// main module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet@gmail.com>
//

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
	var config string
	flag.StringVar(&config, "config", "", "config file")
	var wmstatsFile string
	flag.StringVar(&wmstatsFile, "wmstatsFile", "", "wmstats file")
	flag.Parse()
	if version {
		fmt.Println("wmstats version:", info())
		return
	}
	if wmstatsFile != "" {
		wmgr := NewWMStatsManager(wmstatsFile)
		wmgr.update()
		cmap, smap, rmap, amap := wmstats(wmgr, 1)
		fmt.Println("campaign map", len(cmap))
		fmt.Println("site map", len(smap))
		fmt.Println("release map", len(rmap))
		fmt.Println("agent map", len(amap))
	} else {
		Server(config)
	}
}
