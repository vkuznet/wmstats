package main

// wmstats presentation module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import "fmt"

// WMStatsMap defines interface to represent different WMStats maps
type WMStatsMap interface {
	HTMLTable() string // create proper html table for our map
}

// SiteStatsMap
type SiteStatsMap map[string]SiteStats

// HTMLTable implements WMStatsMap interface
func (wmap SiteStatsMap) HTMLTable() string {
	t := `<table class="is-striped is-bordered" id="site-stats"><tr>`
	t += `<th onclick="sortTable('site-stats', 0)">Site</th>`
	t += `<th onclick="sortTable('site-stats', 1)">Requests</th>`
	t += `<th onclick="sortTable('site-stats', 2)">Pending</th>`
	t += `<th onclick="sortTable('site-stats', 3)">Running</th>`
	t += `<th onclick="sortTable('site-stats', 4)">CoolOff</th>`
	t += `<th onclick="sortTable('site-stats', 5)">Failure Rate</th>`
	t += "</tr>\n"
	for key, data := range wmap {
		t += "<tr>"
		t += fmt.Sprintf("<td>%v</td>", key)
		t += fmt.Sprintf("<td>%v</td>", data.Requests)
		t += fmt.Sprintf("<td>%v</td>", data.Pending)
		t += fmt.Sprintf("<td>%v</td>", data.Running)
		t += fmt.Sprintf("<td>%v</td>", data.CoolOff)
		t += fmt.Sprintf("<td>%v</td>", data.FailureRate)
		t += "</tr>\n"
	}
	t += "</table>"
	return t
}

// CliTable implements WMStatsMap interface
func (wmap SiteStatsMap) CliTable() ([]string, [][]string, []int) {
	headers := []string{
		"Site", "Requests", "Pending", "Running", "CoolOff", "Failure Rate",
	}
	paddings := make([]int, len(headers))
	for k, v := range headers {
		paddings[k] = len(v)
	}
	var allValues [][]string
	for site, data := range wmap {
		if len(site) > paddings[0] {
			paddings[0] = len(site)
		}
		requests := fmt.Sprintf("%v", data.Requests)
		if len(requests) > paddings[1] {
			paddings[1] = len(requests)
		}
		pending := fmt.Sprintf("%v", data.Pending)
		if len(pending) > paddings[2] {
			paddings[2] = len(pending)
		}
		running := fmt.Sprintf("%v", data.Running)
		if len(running) > paddings[3] {
			paddings[3] = len(running)
		}
		cooloff := fmt.Sprintf("%v", data.CoolOff)
		if len(cooloff) > paddings[4] {
			paddings[4] = len(cooloff)
		}
		failureRate := fmt.Sprintf("%v", data.FailureRate)
		if len(failureRate) > paddings[5] {
			paddings[5] = len(failureRate)
		}
		values := []string{
			site, requests, pending, running, cooloff, failureRate,
		}
		allValues = append(allValues, values)
	}
	return headers, allValues, paddings
}

// CampaignStatsMap
type CampaignStatsMap map[string]CampaignStats

// HTMLTable implements WMStatsMap interface
func (wmap CampaignStatsMap) HTMLTable() string {
	t := `<table class="is-striped is-bordered" id="campaign-stats"><tr>`
	t += `<th onclick="sortTable('campaign-stats', 0)">Campaign</th>`
	t += `<th onclick="sortTable('campaign-stats', 1)">Requests</th>`
	t += `<th onclick="sortTable('campaign-stats', 2)">Job Progress</th>`
	t += `<th onclick="sortTable('campaign-stats', 3)">Event Progress</th>`
	t += `<th onclick="sortTable('campaign-stats', 4)">Lumi Progress</th>`
	t += `<th onclick="sortTable('campaign-stats', 5)">Failure Rate</th>`
	t += `<th onclick="sortTable('campaign-stats', 6)">Cool off</th>`
	t += "</tr>\n"
	for key, data := range wmap {
		t += "<tr>"
		link := fmt.Sprintf("%s/workflows?campaign=%s", Config.Base, key)
		ahref := fmt.Sprintf("<a href=\"%s\">%s</a>", link, key)
		t += fmt.Sprintf("<td>%v</td>", ahref)
		t += fmt.Sprintf("<td>%v</td>", data.Requests)
		t += fmt.Sprintf("<td>%v</td>", data.JobProgress)
		t += fmt.Sprintf("<td>%v</td>", data.EventProgress)
		t += fmt.Sprintf("<td>%v</td>", data.LumiProgress)
		t += fmt.Sprintf("<td>%v</td>", data.FailureRate)
		t += fmt.Sprintf("<td>%v</td>", data.CoolOff)
		t += "</tr>\n"
	}
	t += "</table>"
	return t
}

// CliTable implements WMStatsMap interface
func (wmap CampaignStatsMap) CliTable() ([]string, [][]string, []int) {
	headers := []string{
		"Campaign", "Requests", "JobProgress", "EventProgress", "LumiProgress", "Failure Rate", "CoolOff",
	}
	paddings := make([]int, len(headers))
	for k, v := range headers {
		paddings[k] = len(v)
	}
	var allValues [][]string
	for campaign, data := range wmap {
		if len(campaign) > paddings[0] {
			paddings[0] = len(campaign)
		}
		requests := fmt.Sprintf("%v", data.Requests)
		if len(requests) > paddings[1] {
			paddings[1] = len(requests)
		}
		jobProgress := fmt.Sprintf("%v", data.JobProgress)
		if len(jobProgress) > paddings[2] {
			paddings[2] = len(jobProgress)
		}
		eventProgress := fmt.Sprintf("%v", data.EventProgress)
		if len(eventProgress) > paddings[3] {
			paddings[3] = len(eventProgress)
		}
		lumiProgress := fmt.Sprintf("%v", data.LumiProgress)
		if len(lumiProgress) > paddings[4] {
			paddings[4] = len(lumiProgress)
		}
		failureRate := fmt.Sprintf("%v", data.FailureRate)
		if len(failureRate) > paddings[5] {
			paddings[5] = len(failureRate)
		}
		cooloff := fmt.Sprintf("%v", data.CoolOff)
		if len(cooloff) > paddings[6] {
			paddings[6] = len(cooloff)
		}
		values := []string{
			campaign, requests, jobProgress, eventProgress, lumiProgress, failureRate, cooloff,
		}
		allValues = append(allValues, values)
	}
	return headers, allValues, paddings
}

// AgentStatsMap
type AgentStatsMap map[string]AgentStats

// HTMLTable implements WMStatsMap interface
func (wmap AgentStatsMap) HTMLTable() string {
	t := `<table class="is-striped is-bordered" id="agent-stats"><tr>`
	t += `<th onclick="sortTable('agent-stats', 0)">Agent</th>`
	t += `<th onclick="sortTable('agent-stats', 1)">Requests</th>`
	t += `<th onclick="sortTable('agent-stats', 2)">Job Progress</th>`
	t += `<th onclick="sortTable('agent-stats', 3)">Failure Rate</th>`
	t += `<th onclick="sortTable('agent-stats', 4)">Cool off</th>`
	t += "</tr>\n"
	for key, data := range wmap {
		t += "<tr>"
		t += fmt.Sprintf("<td>%v</td>", key)
		t += fmt.Sprintf("<td>%v</td>", data.Requests)
		t += fmt.Sprintf("<td>%v</td>", data.JobProgress)
		t += fmt.Sprintf("<td>%v</td>", data.FailureRate)
		t += fmt.Sprintf("<td>%v</td>", data.CoolOff)
		t += "</tr>\n"
	}
	t += "</table>"
	return t
}

// CliTable implements WMStatsMap interface
func (wmap AgentStatsMap) CliTable() ([]string, [][]string, []int) {
	headers := []string{
		"Agent", "Requests", "JobProgress", "Failure Rate", "CoolOff",
	}
	paddings := make([]int, len(headers))
	for k, v := range headers {
		paddings[k] = len(v)
	}
	var allValues [][]string
	for agent, data := range wmap {
		if len(agent) > paddings[0] {
			paddings[0] = len(agent)
		}
		requests := fmt.Sprintf("%v", data.Requests)
		if len(requests) > paddings[1] {
			paddings[1] = len(requests)
		}
		jobProgress := fmt.Sprintf("%v", data.JobProgress)
		if len(jobProgress) > paddings[2] {
			paddings[2] = len(jobProgress)
		}
		failureRate := fmt.Sprintf("%v", data.FailureRate)
		if len(failureRate) > paddings[3] {
			paddings[3] = len(failureRate)
		}
		cooloff := fmt.Sprintf("%v", data.CoolOff)
		if len(cooloff) > paddings[4] {
			paddings[4] = len(cooloff)
		}
		values := []string{
			agent, requests, jobProgress, failureRate, cooloff,
		}
		allValues = append(allValues, values)
	}
	return headers, allValues, paddings
}

// CMSSWStatsMap
type CMSSWStatsMap map[string]CMSSWStats

// HTMLTable implements WMStatsMap interface
func (wmap CMSSWStatsMap) HTMLTable() string {
	t := `<table class="is-striped is-bordered" id="cmssw-stats"><tr>`
	t += `<th onclick="sortTable('cmssw-stats', 0)">CMSSW</th>`
	t += `<th onclick="sortTable('cmssw-stats', 1)">Requests</th>`
	t += `<th onclick="sortTable('cmssw-stats', 2)">Job Progress</th>`
	t += `<th onclick="sortTable('cmssw-stats', 3)">Event Progress</th>`
	t += `<th onclick="sortTable('cmssw-stats', 4)">Lumi Progress</th>`
	t += `<th onclick="sortTable('cmssw-stats', 5)">Failure Rate</th>`
	t += `<th onclick="sortTable('cmssw-stats', 6)">Cool off</th>`
	t += "</tr>\n"
	for key, data := range wmap {
		t += "<tr>"
		t += fmt.Sprintf("<td>%v</td>", key)
		t += fmt.Sprintf("<td>%v</td>", data.Requests)
		t += fmt.Sprintf("<td>%v</td>", data.JobProgress)
		t += fmt.Sprintf("<td>%v</td>", data.EventProgress)
		t += fmt.Sprintf("<td>%v</td>", data.LumiProgress)
		t += fmt.Sprintf("<td>%v</td>", data.FailureRate)
		t += fmt.Sprintf("<td>%v</td>", data.CoolOff)
		t += "</tr>\n"
	}
	t += "</table>"
	return t
}

// CliTable implements WMStatsMap interface
func (wmap CMSSWStatsMap) CliTable() ([]string, [][]string, []int) {
	headers := []string{
		"CMSSW", "Requests", "JobProgress", "EventProgress", "LumiProgress", "Failure Rate", "CoolOff",
	}
	paddings := make([]int, len(headers))
	for k, v := range headers {
		paddings[k] = len(v)
	}
	var allValues [][]string
	for cmssw, data := range wmap {
		if len(cmssw) > paddings[0] {
			paddings[0] = len(cmssw)
		}
		requests := fmt.Sprintf("%v", data.Requests)
		if len(requests) > paddings[1] {
			paddings[1] = len(requests)
		}
		jobProgress := fmt.Sprintf("%v", data.JobProgress)
		if len(jobProgress) > paddings[2] {
			paddings[2] = len(jobProgress)
		}
		eventProgress := fmt.Sprintf("%v", data.EventProgress)
		if len(eventProgress) > paddings[3] {
			paddings[3] = len(eventProgress)
		}
		lumiProgress := fmt.Sprintf("%v", data.LumiProgress)
		if len(lumiProgress) > paddings[4] {
			paddings[4] = len(lumiProgress)
		}
		failureRate := fmt.Sprintf("%v", data.FailureRate)
		if len(failureRate) > paddings[5] {
			paddings[5] = len(failureRate)
		}
		cooloff := fmt.Sprintf("%v", data.CoolOff)
		if len(cooloff) > paddings[6] {
			paddings[6] = len(cooloff)
		}
		values := []string{
			cmssw, requests, jobProgress, eventProgress, lumiProgress, failureRate, cooloff,
		}
		allValues = append(allValues, values)
	}
	return headers, allValues, paddings
}

// WorkflowMap provides list of workflows for a given key, e.g. campaign
type WorkflowMap map[string][]Workflow

// HTMLTable implements WMStatsMap interface
func workflowHTMLTable(workflows []Workflow) string {
	t := `<table class="is-striped is-bordered" id="wmap"><tr>`
	t += `<th onclick="sortTable('wmap', 0)">Workflow</th>`
	t += `<th onclick="sortTable('wmap', 1)">Status</th>`
	t += `<th onclick="sortTable('wmap', 2)">Type</th>`
	t += `<th onclick="sortTable('wmap', 3)">Priority</th>`
	t += `<th onclick="sortTable('wmap', 4)">Queue Injection</th>`
	t += `<th onclick="sortTable('wmap', 5)">Job Progress</th>`
	t += `<th onclick="sortTable('wmap', 6)">Event Progress</th>`
	t += `<th onclick="sortTable('wmap', 7)">Lumi Progress</th>`
	t += `<th onclick="sortTable('wmap', 8)">Failure Rate</th>`
	t += `<th onclick="sortTable('wmap', 9)">Estimated completion</th>`
	t += `<th onclick="sortTable('wmap', 10)">Cool off</th>`
	t += "</tr>\n"
	for _, data := range workflows {
		t += "<tr>"
		link := fmt.Sprintf("https://cmsweb.cern.ch/reqmgr2/fetch?rid=%s", data.Workflow)
		ahref := fmt.Sprintf("<a href=\"%s\">%s</a>", link, data.Workflow)
		t += fmt.Sprintf("<td>%v</td>", ahref)
		t += fmt.Sprintf("<td>%v</td>", data.Status)
		t += fmt.Sprintf("<td>%v</td>", data.Type)
		t += fmt.Sprintf("<td>%v</td>", data.Priority)
		t += fmt.Sprintf("<td>%v</td>", data.QueueInjection)
		t += fmt.Sprintf("<td>%v</td>", data.JobProgress)
		t += fmt.Sprintf("<td>%v</td>", data.EventProgress)
		t += fmt.Sprintf("<td>%v</td>", data.LumiProgress)
		t += fmt.Sprintf("<td>%v</td>", data.FailureRate)
		t += fmt.Sprintf("<td>%v</td>", data.EstimatedCompletion)
		t += fmt.Sprintf("<td>%v</td>", data.CoolOff)
		t += "</tr>\n"
	}
	t += "</table>"
	return t
}
