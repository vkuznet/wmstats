package main

// cache module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//

import (
	"io"
	"log"
	"os"
	"time"
)

// WMStatsManager manages wmstats data
type WMStatsManager struct {
	URI           string // wmstats URI (URL or file name)
	Data          []byte // wmstats data
	TTL           int64  // time-to-live of current cache snapshot
	RenewInterval int64  // renew interval for cache
	CampaignMap   CampaignStatsMap
	SiteMap       SiteStatsMap
	CMSSWMap      CMSSWStatsMap
	AgentMap      AgentStatsMap
}

// helper function to update cache
func (w *WMStatsManager) update() {
	if w.TTL < time.Now().Unix() {
		var data []byte
		var err error
		if _, err := os.Stat(w.URI); err == nil {
			data, err = readFile(w.URI)
		} else {
			data, err = fetch(w.URI)
		}
		if err == nil {
			w.Data = data
		}
		log.Println("update WMStats cache with", w.URI, err)
		w.TTL = time.Now().Unix() + w.RenewInterval
	}
}

// NewWMStatsManager method properly initialize WMStatsManager
func NewWMStatsManager(uri string, renew ...int64) *WMStatsManager {
	wmstats := &WMStatsManager{URI: uri, RenewInterval: 300} // by default renew cache every 5 minutes
	if len(renew) > 0 {
		wmstats.RenewInterval = renew[0]
	}
	return wmstats
}

// helper function to read data from a file
func readFile(fname string) ([]byte, error) {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	return data, err
}
