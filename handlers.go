package main

// handlers.go - provides handlers examples for wmstats server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// HTTPError represents HTTP error structure
type HTTPError struct {
	Method         string `json:"method"`           // HTTP method
	HTTPCode       int    `json:"code"`             // HTTP status code from IANA
	Timestamp      string `json:"timestamp"`        // timestamp of the error
	Path           string `json:"path"`             // URL path
	UserAgent      string `json:"user_agent"`       // http user-agent field
	XForwardedHost string `json:"x_forwarded_host"` // http.Request X-Forwarded-Host
	XForwardedFor  string `json:"x_forwarded_for"`  // http.Request X-Forwarded-For
	RemoteAddr     string `json:"remote_addr"`      // http.Request remote address
}

// ServerError represents HTTP server error structure
type ServerError struct {
	Error     error     `json:"error"`     // error
	HTTPError HTTPError `json:"http"`      // HTTP section of the error
	Exception int       `json:"exception"` // for compatibility with Python server
	Type      string    `json:"type"`      // for compatibility with Python server
	Message   string    `json:"message"`   // for compatibility with Python server
}

// helper function to parse given template and return HTML page
func tmplPage(tmpl string, tmplData TmplRecord) string {
	if tmplData == nil {
		tmplData = make(TmplRecord)
	}
	var templates Templates
	page := templates.Tmpl(Config.Templates, tmpl, tmplData)
	return page
}

// MetricsHandler provides metrics
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(promMetrics(Config.MetricsPrefix)))
	return
}

// MainHandler provides access to main page of server
func MainHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	stats := query.Get("stats")

	// get data
	wMgr.update()
	cmap, smap, rmap, amap := wmstats(wMgr)
	var table string
	if stats == "agent" {
		table = amap.HTMLTable()
	} else if stats == "site" {
		table = smap.HTMLTable()
	} else if stats == "cmssw" {
		table = rmap.HTMLTable()
	} else if stats == "campaign" {
		table = cmap.HTMLTable()
	} else {
		table = cmap.HTMLTable()
	}

	// create temaplate
	tmpl := make(TmplRecord)
	tmpl["Base"] = Config.Base
	tmpl["ServerInfo"] = ServerInfo
	tmpl["Table"] = template.HTML(table)

	page := tmplPage("main.tmpl", tmpl)
	w.Write([]byte(page))
}

// StatusHandler provides basic functionality of status response
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	//     records = append(records, rec)
	//     data, err := json.Marshal(records)
	//     if err != nil {
	//         log.Fatalf("Fail to marshal records, %v", err)
	//     }
	data := []byte("ok")
	w.Write(data)
}

// ServerInfoHandler provides basic functionality of status response
func ServerInfoHandler(w http.ResponseWriter, r *http.Request) {
	rec := make(map[string]string)
	rec["server"] = ServerInfo
	data, err := json.Marshal(rec)
	if err != nil {
		log.Fatalf("Fail to marshal records, %v", err)
	}
	w.Write(data)
}
