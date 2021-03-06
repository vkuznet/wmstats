package main

// wmstats cli module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"fmt"
	"strings"
)

// cli provides CLI interface to wmstats
func cli(wmstatsFile string, filters WMStatsFilters, stats string, verbose int) {
	wmgr := NewWMStatsManager(wmstatsFile)
	wmgr.update()
	_wmstatsInfo := wmstats(wmgr, filters, verbose)
	var headers []string
	var values [][]string
	var paddings []int
	if stats == "agent" {
		headers, values, paddings = _wmstatsInfo.AgentStatsMap.CliTable()
	} else if stats == "site" {
		headers, values, paddings = _wmstatsInfo.SiteStatsMap.CliTable()
	} else if stats == "cmssw" {
		headers, values, paddings = _wmstatsInfo.CMSSWStatsMap.CliTable()
	} else if stats == "campaign" {
		headers, values, paddings = _wmstatsInfo.CampaignStatsMap.CliTable()
	} else {
		headers, values, paddings = _wmstatsInfo.CampaignStatsMap.CliTable()
	}

	// print headers
	var rowValues []string
	for k, v := range headers {
		if len(v) <= paddings[k] {
			vals := make([]string, paddings[k]+1-len(v))
			rowValues = append(rowValues, v+strings.Join(vals, " "))
		} else {
			rowValues = append(rowValues, v)
		}
	}
	fmt.Println(strings.Join(rowValues, " "))
	// print values
	for _, vals := range values {
		var rowValues []string
		for k, v := range vals {
			if len(v) <= paddings[k] {
				vals := make([]string, paddings[k]+1-len(v))
				rowValues = append(rowValues, v+strings.Join(vals, " "))
			} else {
				rowValues = append(rowValues, v)
			}
		}
		fmt.Println(strings.Join(rowValues, " "))
	}
}
