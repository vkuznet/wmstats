package main

import "strings"

type WMStatsFilters map[string]string

// helper function to get wmstats filters out of HTTP request
func wmstatsFilters(values string) WMStatsFilters {
	filters := make(WMStatsFilters)
	arr := strings.Split(values, ",")
	for _, val := range arr {
		pair := strings.Split(val, "=")
		if len(pair) != 2 {
			//             log.Printf("Unable to extract key=value pair from wmstats filters: '%s'\n", wmstatsFilters)
			continue
		}
		key := pair[0]
		value := pair[1]
		filters[key] = value
	}
	return filters
}
