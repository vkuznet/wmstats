package main

// handlers.go - provides handlers examples for wmstats server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
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

// global pointer to wmstats info
var _wmstatsInfo *WMStatsInfo

// MainHandler provides access to main page of server
func MainHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	stats := query.Get("stats")

	// get data
	wMgr.update()
	if _wmstatsInfo == nil || wMgr.TTL < time.Now().Unix() {
		_wmstatsInfo = wmstats(wMgr, 0)
	}
	var table string
	if stats == "agent" {
		table = _wmstatsInfo.AgentStatsMap.HTMLTable()
	} else if stats == "site" {
		table = _wmstatsInfo.SiteStatsMap.HTMLTable()
	} else if stats == "cmssw" {
		table = _wmstatsInfo.CMSSWStatsMap.HTMLTable()
	} else if stats == "campaign" {
		table = _wmstatsInfo.CampaignStatsMap.HTMLTable()
	} else {
		table = _wmstatsInfo.CampaignStatsMap.HTMLTable()
	}

	// create temaplate
	tmpl := make(TmplRecord)
	tmpl["Base"] = Config.Base
	tmpl["ServerInfo"] = ServerInfo
	tmpl["Table"] = template.HTML(table)
	tmpl["Menu"] = template.HTML(tmplPage("menu.tmpl", tmpl))
//     tmpl["Search"] = template.HTML(tmplPage("search.tmpl", tmpl))
	tmpl["Filter"] = template.HTML(tmplPage("filters.tmpl", tmpl))
	tmpl["Header"] = _header
	tmpl["Footer"] = _footer

	page := tmplPage("main.tmpl", tmpl)
	w.Write([]byte(string(_top) + page + string(_bottom)))
}

// AlertsHandler provides access to alerts page of server
func AlertsHandler(w http.ResponseWriter, r *http.Request) {
	// create temaplate
	tmpl := make(TmplRecord)
	tmpl["Base"] = Config.Base
	tmpl["ServerInfo"] = ServerInfo
	tmpl["Menu"] = template.HTML(tmplPage("menu.tmpl", tmpl))
	tmpl["Header"] = _header
	tmpl["Footer"] = _footer

	page := tmplPage("alerts.tmpl", tmpl)
	w.Write([]byte(string(_top) + page + string(_bottom)))
}

// AgentsHandler provides access to agents page of server
func AgentsHandler(w http.ResponseWriter, r *http.Request) {
	// create temaplate
	tmpl := make(TmplRecord)
	tmpl["Base"] = Config.Base
	tmpl["ServerInfo"] = ServerInfo
	tmpl["Menu"] = template.HTML(tmplPage("menu.tmpl", tmpl))
	tmpl["Header"] = _header
	tmpl["Footer"] = _footer

	page := tmplPage("agents.tmpl", tmpl)
	w.Write([]byte(string(_top) + page + string(_bottom)))
}

// ErrorLogsHandler provides access to error logs page of server
func ErrorLogsHandler(w http.ResponseWriter, r *http.Request) {
	// create temaplate
	tmpl := make(TmplRecord)
	tmpl["Base"] = Config.Base
	tmpl["ServerInfo"] = ServerInfo
	tmpl["Menu"] = template.HTML(tmplPage("menu.tmpl", tmpl))
	tmpl["Header"] = _header
	tmpl["Footer"] = _footer

	page := tmplPage("errorlogs.tmpl", tmpl)
	w.Write([]byte(string(_top) + page + string(_bottom)))
}

// WorkflowsHandler provides access to workflows page of server
func WorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	campaign := query.Get("campaign")
	site := query.Get("site")
	cmssw := query.Get("cmssw")
	agent := query.Get("agent")

	// get data
	wMgr.update()
	if _wmstatsInfo == nil || wMgr.TTL < time.Now().Unix() {
		_wmstatsInfo = wmstats(wMgr, 0)
	}
	table := "Unkown key"
	title := ""
	if campaign != "" {
		if workflows, ok := _wmstatsInfo.CampaignWorkflows[campaign]; ok {
			table = workflowHTMLTable(workflows)
			val := fmt.Sprintf("<span class=\"alert is-focus\">%s</span>", campaign)
			title = fmt.Sprintf("<h4>Workflows associated with %s campaign</h4>", val)
		}
	} else if site != "" {
		if workflows, ok := _wmstatsInfo.SiteWorkflows[site]; ok {
			table = workflowHTMLTable(workflows)
			val := fmt.Sprintf("<span class=\"alert is-focus\">%s</span>", site)
			title = fmt.Sprintf("<h4>Workflows associated with %s site</h4>", val)
		}
	} else if cmssw != "" {
		if workflows, ok := _wmstatsInfo.CMSSWWorkflows[cmssw]; ok {
			table = workflowHTMLTable(workflows)
			val := fmt.Sprintf("<span class=\"alert is-focus\">%s</span>", cmssw)
			title = fmt.Sprintf("<h4>Workflows associated with with %s</h4>", val)
		}
	} else if agent != "" {
		if workflows, ok := _wmstatsInfo.AgentWorkflows[agent]; ok {
			table = workflowHTMLTable(workflows)
			val := fmt.Sprintf("<span class=\"alert is-focus\">%s</span>", agent)
			title = fmt.Sprintf("<h4>Workflows associated with %s agent</h4>", val)
		}
	}

	// create temaplate
	tmpl := make(TmplRecord)
	tmpl["Base"] = Config.Base
	tmpl["ServerInfo"] = ServerInfo
	tmpl["Menu"] = template.HTML(tmplPage("menu.tmpl", tmpl))
	tmpl["Header"] = _header
	tmpl["Footer"] = _footer
	tmpl["Title"] = template.HTML(title)
	tmpl["Table"] = template.HTML(table)

	page := tmplPage("main.tmpl", tmpl)
	w.Write([]byte(string(_top) + page + string(_bottom)))
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
