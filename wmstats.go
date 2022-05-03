package main

// wmstats module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"encoding/json"
	"fmt"
	"log"
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
	t := `<table id="site-stats"><tr>`
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

// CampaignStatsMap
type CampaignStatsMap map[string]CampaignStats

// HTMLTable implements WMStatsMap interface
func (wmap CampaignStatsMap) HTMLTable() string {
	t := `<table id="campaign-stats"><tr>`
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

// AgentStatsMap
type AgentStatsMap map[string]AgentStats

// HTMLTable implements WMStatsMap interface
func (wmap AgentStatsMap) HTMLTable() string {
	t := `<table id="agent-stats"><tr>`
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

// CMSSWStatsMap
type CMSSWStatsMap map[string]CMSSWStats

// HTMLTable implements WMStatsMap interface
func (wmap CMSSWStatsMap) HTMLTable() string {
	t := `<table id="cmssw-stats"><tr>`
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

// wmstats provide aggregated statistics
func wmstats(wmgr *WMStatsManager, verbose int) (CampaignStatsMap, SiteStatsMap, CMSSWStatsMap, AgentStatsMap) {
	time0 := time.Now()
	// update our cache
	wmgr.update()
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

	// main loop
	for _, info := range data {
		for workflow, rdict := range info {
			if verbose > 0 {
				fmt.Println(workflow)
				//             fmt.Printf("%+v\n", rdict)
			}
			cmssw := rdict.CMSSWVersion
			workflow := rdict.RequestName
			//             totalEvents := rdict.TotalInputEvents
			//             totalLumis := rdict.TotalInputLumis

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
	if verbose > 0 {
		fmt.Println("### Total site stats", len(smap))
	}
	for site, stats := range smap {
		if verbose > 0 {
			fmt.Println("site", site)
		}
		workflows, _ := sWorkflows[site]
		stats.Requests = workflows.Size()
		totJobs := stats.SuccessJobs + stats.FailJobs
		if totJobs != 0 {
			stats.FailureRate = 100 * float64(stats.FailJobs) / float64(totJobs)
		}
		if verbose > 0 {
			fmt.Printf("%+v\n", stats)
		}
	}

	// collect campaign summary from workflow map
	for _, winfo := range wmap {
		campaign := winfo.Campaign
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
	if verbose > 0 {
		fmt.Println("### Total campaign stats", len(cmap))
	}

	//     fmt.Println("### campaign summary")
	//     for c, data := range campaignSummary {
	//         log.Println(c, data)
	//     }

	// prepare campaign stats dict
	for campaign, stats := range cmap {
		if verbose > 0 {
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
		if verbose > 0 {
			fmt.Printf("%+v\n", stats)
		}
	}
	if verbose > 0 {
		fmt.Println("### agent summary", len(agentSummary))
	}
	for agent, data := range agentSummary {
		if verbose > 0 {
			fmt.Println("agent:", agent)
			fmt.Printf("%+v\n", data)
		}
	}

	if verbose > 0 {
		fmt.Println("### cmssw summary", len(cmsswSummary))
	}
	for cmssw, data := range cmsswSummary {
		if verbose > 0 {
			fmt.Println("cmssw:", cmssw)
			fmt.Printf("%+v\n", data)
		}
	}
	fmt.Println("### Total number of workflows", len(wmap), "in", time.Since(time0))
	return cmap, smap, rmap, amap
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
