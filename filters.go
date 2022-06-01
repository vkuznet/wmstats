package main

// wmstats filters module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"fmt"
	"strings"
)

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

// helper function to convert list of filters to HTML format
func filtersToHTML(filters WMStatsFilters) string {
	var s string
	for k, v := range filters {
		// TODO: make close icon to discard the filter
		// so the action should turn off the span
		f := fmt.Sprintf("<span class=\"alert is-focus\">%s=%s</span>", k, v)
		s += fmt.Sprintf("%s<br/>", f)
	}
	return s
}
