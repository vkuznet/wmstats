package main

// wmstats module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	// data set
	"github.com/fatih/set"
)

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
		"Site", "Requests", "Pending", "Running", "CoolOff", "Failure Raute",
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
		"Campaign", "Requests", "JobProgress", "EventProgress", "LumiProgress", "Failure Raute", "CoolOff",
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
		"CMSSW", "Requests", "JobProgress", "EventProgress", "LumiProgress", "Failure Raute", "CoolOff",
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

// WMStatsInfo represent wmstats info structure
type WMStatsInfo struct {
	CampaignStatsMap  CampaignStatsMap
	SiteStatsMap      SiteStatsMap
	CMSSWStatsMap     CMSSWStatsMap
	AgentStatsMap     AgentStatsMap
	CampaignWorkflows WorkflowMap
	SiteWorkflows     WorkflowMap
	CMSSWWorkflows    WorkflowMap
	AgentWorkflows    WorkflowMap
}

// wmstats provide aggregated statistics
func wmstats(wmgr *WMStatsManager, filters WMStatsFilters, verbose int) *WMStatsInfo {
	time0 := time.Now()
	// update our cacheAgentStatsMawmgr.update()
	var wmstats WMStatsResults
	err := json.Unmarshal(wmgr.Data, &wmstats)
	if err != nil {
		log.Fatal(err)
	}
	data := wmstats.Result

	//     data := readData(fname)
	// create our stats maps
	cmap := make(CampaignStatsMap)
	smap := make(SiteStatsMap)
	amap := make(AgentStatsMap)
	rmap := make(CMSSWStatsMap)

	// create aux maps
	agentSummary := make(map[string]AgentSummary)
	cmsswSummary := make(map[string]CMSSWSummary)
	campaignSummary := make(map[string]CampaignSummary)
	sWorkflows := make(map[string]set.Interface)
	wmap := make(map[string]WorkflowInfo)

	// create workflow maps
	campaignMap := make(WorkflowMap)
	siteMap := make(WorkflowMap)
	cmsswMap := make(WorkflowMap)
	agentMap := make(WorkflowMap)

	// main loop
	for _, info := range data {
		for workflow, rdict := range info {
			if verbose > 1 {
				fmt.Println(workflow)
				//             fmt.Printf("%+v\n", rdict)
			}
			cmssw := rdict.CMSSWVersion
			workflow := rdict.RequestName
			campaign := rdict.Campaign
			//             totalEvents := rdict.TotalInputEvents
			//             totalLumis := rdict.TotalInputLumis

			wObj := Workflow{
				Workflow:            workflow,
				QueueInjection:      0.0,
				JobProgress:         0.0,
				EventProgress:       0.0,
				LumiProgress:        0.0,
				FailureRate:         0.0,
				Status:              rdict.RequestStatus,
				Type:                rdict.RequestType,
				EstimatedCompletion: "N/A",
				Priority:            rdict.RequestPriority,
				CoolOff:             0,
			}
			updateMap(campaignMap, campaign, wObj)
			updateMap(cmsswMap, cmssw, wObj)

			// collect workflow information
			wInfo := WorkflowInfo{
				Name:     rdict.RequestName,
				Campaign: rdict.Campaign,
				Type:     rdict.RequestType,
				Priority: rdict.RequestPriority,
				Sites:    rdict.Sites,
			}
			// keey workflow info regardless of AgentJobInfoMap which may be missing
			wmap[workflow] = wInfo
			// setup initial values for cmssw summary
			if cs, ok := cmsswSummary[cmssw]; ok {
				cs.Requests += 1
				cmsswSummary[cmssw] = cs
			} else {
				cmsswSummary[cmssw] = CMSSWSummary{Requests: 1}
			}

			// collect site statistics from AgentJobInfo map
			var agents []string
			for agent, ainfo := range rdict.AgentJobInfoMap {
				updateMap(agentMap, agent, wObj)
				agents = append(agents, agent)
				workflow := ainfo.Workflow
				status := ainfo.Status
				// update workflow info
				if winfo, ok := wmap[workflow]; ok {
					winfo.Status.Update(status)
					wmap[workflow] = winfo
				} else {
					wInfo.Status = status
					wmap[workflow] = wInfo
				}

				// update cmssw info
				updateReleaseSummary(cmssw, cmsswSummary, status)

				// update agent info
				updateAgentSummary(agent, agentSummary, status)

				// update campaing info

				// update site info
				for site, status := range ainfo.Sites {
					updateMap(siteMap, site, wObj)
					coolOff := status.CoolOff.Sum()
					pending := status.Submitted.Pending
					running := status.Submitted.Running
					failure := status.Failure.Sum()
					success := status.Success
					if stats, ok := smap[site]; ok {
						stats.CoolOff += coolOff
						stats.Pending += pending
						stats.Running += running
						stats.FailJobs += failure
						stats.SuccessJobs += success
						smap[site] = stats
					} else {
						stats := SiteStats{
							CoolOff:     coolOff,
							Pending:     pending,
							Running:     running,
							FailJobs:    failure,
							SuccessJobs: success,
						}
						smap[site] = stats
					}
					if workflows, ok := sWorkflows[site]; ok {
						workflows.Add(workflow)
						sWorkflows[site] = workflows
					} else {
						swSet := set.New(set.ThreadSafe)
						swSet.Add(workflow)
						sWorkflows[site] = swSet
					}
				}
			}
			winfo, _ := wmap[wInfo.Name]
			winfo.Agents = agents
			wmap[wInfo.Name] = winfo
		}
	}
	// prepare site stats dict
	if verbose > 1 {
		fmt.Println("### Total site stats", len(smap))
	}
	var sites []string
	for site, stats := range smap {
		if verbose > 1 {
			fmt.Println("site", site)
		}
		workflows, _ := sWorkflows[site]
		stats.Requests = workflows.Size()
		totJobs := stats.SuccessJobs + stats.FailJobs
		if totJobs != 0 {
			stats.FailureRate = 100 * float64(stats.FailJobs) / float64(totJobs)
		}
		if verbose > 1 {
			fmt.Printf("%+v\n", stats)
		}
		sites = append(sites, site)
	}
	// filter out sites
	for _, site := range sites {
		if pat, ok := filters["site"]; ok {
			if matched, err := regexp.MatchString(pat, site); err == nil && !matched {
				delete(smap, site)
			}
		}
	}

	// collect campaign summary from workflow map
	for _, winfo := range wmap {
		campaign := winfo.Campaign
		// filter out unnecessary campaigns
		if pat, ok := filters["campaign"]; ok {
			if matched, err := regexp.MatchString(pat, campaign); err == nil && !matched {
				continue
			}
		}
		if cs, ok := campaignSummary[campaign]; ok {
			cs.Requests += 1
			cs.Status.Update(winfo.Status)
			campaignSummary[campaign] = cs
		} else {
			summary := CampaignSummary{}
			summary.Status.Update(winfo.Status)
			summary.Requests = 1
			campaignSummary[campaign] = summary
		}
		cmap[campaign] = CampaignStats{}
		//         fmt.Printf("workflow: %s\n", workflow)
		//         fmt.Printf("%+v\n", winfo)
	}
	if verbose > 1 {
		fmt.Println("### Total campaign stats", len(cmap))
	}

	//     fmt.Println("### campaign summary")
	//     for c, data := range campaignSummary {
	//         log.Println(c, data)
	//     }

	// prepare campaign stats dict
	for campaign, stats := range cmap {
		if verbose > 1 {
			fmt.Println("campaign", campaign)
		}
		if cs, ok := campaignSummary[campaign]; ok {
			stats.JobProgress = cs.JobProgress()
			stats.EventProgress = cs.EventProgress()
			stats.LumiProgress = cs.LumiProgress()
			stats.FailureRate = cs.FailureRate()
			stats.Requests = cs.Requests
			stats.CoolOff = cs.Status.CoolOff.Sum()
			cmap[campaign] = stats
		}
		if verbose > 1 {
			fmt.Printf("%+v\n", stats)
		}
	}
	if verbose > 1 {
		fmt.Println("### agent summary", len(agentSummary))
	}
	for agent, data := range agentSummary {
		if verbose > 1 {
			fmt.Println("agent:", agent)
			fmt.Printf("%+v\n", data)
		}
	}

	if verbose > 1 {
		fmt.Println("### cmssw summary", len(cmsswSummary))
	}
	for cmssw, data := range cmsswSummary {
		if verbose > 1 {
			fmt.Println("cmssw:", cmssw)
			fmt.Printf("%+v\n", data)
		}
	}
	fmt.Println("### Total number of workflows", len(wmap), "in", time.Since(time0))
	stats := WMStatsInfo{
		CampaignStatsMap:  cmap,
		SiteStatsMap:      smap,
		CMSSWStatsMap:     rmap,
		AgentStatsMap:     amap,
		CampaignWorkflows: campaignMap,
		SiteWorkflows:     siteMap,
		CMSSWWorkflows:    cmsswMap,
		AgentWorkflows:    agentMap,
	}
	return &stats
}

func updateReleaseSummary(cmssw string, cmsswSummary map[string]CMSSWSummary, status Status) {
	if rinfo, ok := cmsswSummary[cmssw]; ok {
		rinfo.Status.Update(status)
		rinfo.CoolOff += rinfo.Status.CoolOff.Sum()
		cmsswSummary[cmssw] = rinfo
	} else {
		rinfo := CMSSWSummary{}
		rinfo.Status.Update(status)
		rinfo.CoolOff += rinfo.Status.CoolOff.Job
		cmsswSummary[cmssw] = rinfo
	}
}

func updateAgentSummary(agent string, agentSummary map[string]AgentSummary, status Status) {
	if ainfo, ok := agentSummary[agent]; ok {
		ainfo.Requests += 1
		ainfo.Status.Update(status)
		ainfo.CoolOff += ainfo.Status.CoolOff.Sum()
		agentSummary[agent] = ainfo
	} else {
		ainfo := AgentSummary{}
		ainfo.Requests += 1
		ainfo.Status.Update(status)
		ainfo.CoolOff += ainfo.Status.CoolOff.Sum()
		agentSummary[agent] = ainfo
	}
}

func updateMap(wmap WorkflowMap, key string, workflow Workflow) {
	if vals, ok := wmap[key]; ok {
		vals = append(vals, workflow)
		wmap[key] = vals
	} else {
		wmap[key] = []Workflow{workflow}
	}
}
