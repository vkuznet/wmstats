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
