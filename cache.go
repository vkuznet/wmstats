package main

// cache module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//

import (
	"time"
)

// WMStatsManager manages wmstats data
type WMStatsManager struct {
	URL           string // wmstats provider URL
	Data          []byte // wmstats data
	TTL           int64  // time-to-live of current cache snapshot
	RenewInterval int64  // renew interval for cache
}

// helper function to update DNSManager cache
func (w *WMStatsManager) update() {
	if w.TTL < time.Now().Unix() {
		data, err := fetch(w.URL)
		if err == nil {
			w.Data = data
		}
		w.TTL = time.Now().Unix() + w.RenewInterval
	}
}

// NewWMStatsManager method properly initialize WMStatsManager
func NewWMStatsManager(renew ...int64) *WMStatsManager {
	wmstats := &WMStatsManager{RenewInterval: 10} // by default renew cache every 10 seconds
	if len(renew) > 0 {
		wmstats.RenewInterval = renew[0]
	}
	return wmstats
}
