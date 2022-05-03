package main

// wmstats module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"encoding/json"
	"fmt"
	"log"

	// data set
	"github.com/fatih/set"
)

// func readData(fname string) []WMStatsRecords {
//     file, err := os.Open(fname)
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer file.Close()
//     data, err := io.ReadAll(file)
//     if err != nil {
//         log.Fatal(err)
//     }
//     var wmstats WMStatsResults
//     err = json.Unmarshal(data, &wmstats)
//     if err != nil {
//         log.Fatal(err)
//     }
//     return wmstats.Result
// }

type SiteStatsMap map[string]SiteStats
type CampaignStatsMap map[string]CampaignStats
type AgentStatsMap map[string]AgentStats
type CMSSWStatsMap map[string]CMSSWStats

func wmstats(wmgr *WMStatsManager) (CampaignStatsMap, SiteStatsMap, CMSSWStatsMap, AgentStatsMap) {
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
			fmt.Println(workflow)
			//             fmt.Printf("%+v\n", rdict)
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
	fmt.Println("### Total site stats", len(smap))
	for site, stats := range smap {
		fmt.Println("site", site)
		workflows, _ := sWorkflows[site]
		stats.Requests = workflows.Size()
		totJobs := stats.SuccessJobs + stats.FailJobs
		if totJobs != 0 {
			stats.FailureRate = 100 * float64(stats.FailJobs) / float64(totJobs)
		}
		fmt.Printf("%+v\n", stats)
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
	fmt.Println("### Total campaign stats", len(cmap))
	//     fmt.Println("### campaign summary")
	//     for c, data := range campaignSummary {
	//         log.Println(c, data)
	//     }

	// prepare campaign stats dict
	for campaign, stats := range cmap {
		fmt.Println("campaign", campaign)
		if cs, ok := campaignSummary[campaign]; ok {
			stats.JobProgress = cs.JobProgress()
			stats.EventProgress = cs.EventProgress()
			stats.LumiProgress = cs.LumiProgress()
			stats.FailureRate = cs.FailureRate()
			stats.Requests = cs.Requests
			stats.CoolOff = cs.Status.CoolOff.Sum()
		}
		fmt.Printf("%+v\n", stats)
	}
	fmt.Println("### agent summary", len(agentSummary))
	for agent, data := range agentSummary {
		fmt.Println("agent:", agent)
		fmt.Printf("%+v\n", data)
	}

	fmt.Println("### cmssw summary", len(cmsswSummary))
	for cmssw, data := range cmsswSummary {
		fmt.Println("cmssw:", cmssw)
		fmt.Printf("%+v\n", data)
	}
	fmt.Println("### Total number of workflows", len(wmap))
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
