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
		_wmstatsInfo := wmstats(wmgr, 1)
		fmt.Println("campaign map", len(_wmstatsInfo.CampaignStatsMap))
		fmt.Println("site map", len(_wmstatsInfo.SiteStatsMap))
		fmt.Println("release map", len(_wmstatsInfo.CMSSWStatsMap))
		fmt.Println("agent map", len(_wmstatsInfo.AgentStatsMap))
	} else {
		Server(config)
	}
}
