package main

// Some schema of wmstatsserver which runs on
// https://cmsweb.cern.ch/wmstatsserver/data/requestcache
// is defined in
// WMCore/src/python/WMCore/Services/WMStats/WMStatsReader.py

// WMStatsRecords defines map of WMStats records
type WMStatsRecords map[string]WMStats

// WMStatsResults represents response from wmstatsserver
type WMStatsResults struct {
	Result []WMStatsRecords
}

// Failure represents failure structure
type Failure struct {
	Exception int
	Create    int
	Submit    int
}

// Sum shows all failures
func (f *Failure) Sum() int {
	return f.Exception + f.Create + f.Submit
}

// Update provides updates of all failure attributes
func (s *Failure) Update(r Failure) {
	s.Exception += r.Exception
	s.Submit += r.Submit
	s.Create += r.Create
}

// Paused represents Paused structure
type Paused struct {
	Job    int
	Submit int
	Create int
}

// Sum shows sum of Paused attributes
func (p *Paused) Sum() int {
	return p.Job + p.Submit + p.Create
}

// Update provides updates of all Paused attributes
func (s *Paused) Update(r Paused) {
	s.Job += r.Job
	s.Submit += r.Submit
	s.Create += r.Create
}

// CoolOff represents cooloff structure
type CoolOff struct {
	Job    int
	Submit int
	Create int
}

// Sum shows sum of cooloffs attributes
func (c *CoolOff) Sum() int {
	return c.Job + c.Submit + c.Create
}

// Update provides updates of all cooloff attributes
func (s *CoolOff) Update(r CoolOff) {
	s.Job += r.Job
	s.Submit += r.Submit
	s.Create += r.Create
}

// Queued represents queues structure
type Queued struct {
	First int
	Retry int
}

// Sum shows sum of queued attributes
func (q *Queued) Sum() int {
	return q.First + q.Retry
}

// Update provides updates of all queued attributes
func (s *Queued) Update(r Queued) {
	s.First += r.First
	s.Retry += r.Retry
}

// Submitted represents submitted structure
type Submitted struct {
	Running int
	Pending int
	Retry   int
}

// Update provides updates of all submitted attributes
func (s *Submitted) Update(r Submitted) {
	s.Running += r.Running
	s.Pending += r.Pending
	s.Retry += r.Retry
}

// Status represents status structure
type Status struct {
	Failure    Failure   `json:"failure"`
	CoolOff    CoolOff   `json:"cooloff"`
	Queued     Queued    `json:"queued"`
	Submitted  Submitted `json:"submitted"`
	Paused     Paused    `json:"paused"`
	Success    int       `json:"success"`
	Canceled   int       `json:"canceled"`
	InWMBS     int       `json:"inWMBS"`
	InQueue    int       `json:"inQueue"`
	Transition int       `json:"transition"`
}

// Update provides updates of all status attributes
func (s *Status) Update(status Status) {
	s.Failure.Update(status.Failure)
	s.CoolOff.Update(status.CoolOff)
	s.Submitted.Update(status.Submitted)
	s.Queued.Update(status.Queued)
	s.Paused.Update(status.Paused)
	s.Transition += status.Transition
	s.Success += status.Success
	s.Canceled += status.Canceled
	s.InWMBS += status.InWMBS
	s.InQueue += status.InQueue
	s.Transition += status.Transition
}

// AgentJobInfo represents WMAgent job information
type AgentJobInfo struct {
	AgentUrl string `json:"agent_url"`
	Workflow string
	Status   Status
	Sites    map[string]Status
}

// Tasks represents tasks data structure
type Task struct {
	PrepID         string
	Campaign       string
	TaskName       string
	AcquisitionEra string
	PrimaryDataset string
	GlobalTag      string
	CMSSWVersion   string
}

// AgentJobInfoMap represents new data tupe for map of WMAgent job info
type AgentJobInfoMap map[string]AgentJobInfo

// WMStats represents wmstats data-structure we parse
type WMStats struct {
	//     Task1           Task `json:"Task1"`
	//     Task2           Task `json:"Task2"`
	//     Task3           Task `json:"Task3"`
	//     Task4           Task `json:"Task4"`
	//     Task5           Task `json:"Task5"`
	//     Task6           Task `json:"Task6"`
	CMSSWVersion    string
	Campaign        string
	RequestStatus   string          `json:"RequestStatus"`
	RequestPriority float64         `json:"RequestPriority"`
	RequestType     string          `json:"RequestType"`
	RequestName     string          `json:"RequestName"`
	Sites           []string        `json:"SiteWhiteList"`
	AgentJobInfoMap AgentJobInfoMap `json:"AgentJobInfo"`
}

// WorkflowInfo provides useful map between workflow (task) name
// and other attributes such as list of sites, releases, agents, etc.
type WorkflowInfo struct {
	Priority float64
	Name     string
	Type     string
	Campaign string
	Sites    []string
	Agents   []string
	Releases []string
	Status   Status
}

type Workflow struct {
	QueueInjection      float64
	JobProgress         float64
	EventProgress       float64
	LumiProgress        float64
	FailureRate         float64
	Status              string
	Type                string
	EstimatedCompletion string
	Priority            int
	CoolOff             int
}

// SiteStats represents common statistics about sites
// see WMCore/src/couchapps/WMStats/_attachments/js/Views/Tables/T1/WMStats.SiteSummaryTable.js
type SiteStats struct {
	FailureRate float64
	Requests    int
	CoolOff     int
	Pending     int
	Running     int
	FailJobs    int
	SuccessJobs int
}

// CMSSWStats represents common statistics about CMSSW releases
// see WMCore/src/couchapps/WMStats/_attachments/js/Views/Tables/T1/WMStats.CMSSWSummaryTable.js
type CMSSWStats struct {
	JobProgress   float64
	EventProgress float64
	LumiProgress  float64
	FailureRate   float64
	Requests      int
	CoolOff       int
}

// CMSSWSummary keeps information about CMSSW summary
type CMSSWSummary struct {
	Status    Status
	Requests  int
	TotalJobs int
	CoolOff   int
}

// AgentStats represents common statistics about agents
// see WMCore/src/couchapps/WMStats/_attachments/js/Views/Tables/T1/WMStats.AgentRequestSummaryTable.js
type AgentStats struct {
	FailureRate float64
	JobProgress float64
	Requests    int
	CoolOff     int
}

// AgentSummary keeps information about agent summary
type AgentSummary struct {
	Status    Status
	Requests  int
	TotalJobs int
	CoolOff   int
}

// CampaignStats represents common statistics about campaigns
// see WMCore/src/couchapps/WMStats/_attachments/js/Views/Tables/T1/WMStats.CampaignSummaryTable.js
//     WMCore/src/couchapps/WMStats/_attachments/js/DataStruct/T1/WMStats.CampaignSummary.js
//     WMCore/src/couchapps/WMStats/_attachments/js/DataStruct/WMStats.GenericRequests.js
type CampaignStats struct {
	JobProgress   float64
	EventProgress float64
	LumiProgress  float64
	FailureRate   float64
	Requests      int
	CoolOff       int
}

// CampaignSummary keeps information about campaign summary
type CampaignSummary struct {
	Status      Status
	Requests    int
	TotalJobs   int
	TotalEvents int
	TotalLumis  int
}

func (cs *CampaignSummary) WMBSTotalJobs() float64 {
	success := cs.Status.Success
	canceled := cs.Status.Canceled
	transition := cs.Status.Transition
	failure := cs.Status.Failure.Sum()
	cooloff := cs.Status.CoolOff.Sum()
	queued := cs.Status.Queued.Sum()
	paused := cs.Status.Paused.Sum()
	running := cs.Status.Submitted.Running
	pending := cs.Status.Submitted.Pending
	return float64(success + canceled + transition + failure + cooloff + paused + queued + running + pending)
}
func (cs *CampaignSummary) JobProgress() float64 {
	totalJobs := cs.WMBSTotalJobs()
	if totalJobs == 0 {
		totalJobs = 1
	}
	return 100 * float64(cs.Status.Success+cs.Status.Failure.Sum()) / float64(totalJobs)
}
func (cs *CampaignSummary) EventProgress() float64 {
	return 100 * float64(cs.AvgEvents()) / float64(cs.TotalEvents)
}
func (cs *CampaignSummary) LumiProgress() float64 {
	return 100 * float64(cs.AvgLumis()) / float64(cs.TotalLumis)
}
func (cs *CampaignSummary) FailureRate() float64 {
	return 100 * float64(cs.Status.Failure.Sum()) / float64(cs.TotalJobs)
}
func (cs *CampaignSummary) AvgEvents() float64 {
	return 1
}
func (cs *CampaignSummary) AvgLumis() float64 {
	return 1
}

/*
from WMCore/src/couchapps/WMStats/_attachments/js/Views/Tables/T1/WMStats.CampaignSummaryTable.js

job progress
		var totalJobs = campaignSummary.getWMBSTotalJobs() || 1;
		var result = (campaignSummary.getJobStatus("success") + campaignSummary.getTotalFailure()) /totalJobs * 100;


event progress
		var totalEvents = row.summary.summaryStruct.totalEvents || 1;
		var result = row.summary.getAvgEvents() / totalEvents * 100;
lumi progress
		var totalLumis = row.summary.summaryStruct.totalLumis || 1;
		var result = row.summary.getAvgLumis() / totalLumis * 100;
failure rate
	   var totalFailure = campaignSummary.getTotalFailure();
	   var totalJobs = (campaignSummary.getJobStatus("success") + totalFailure) || 1;
	   var result = totalFailure / totalJobs * 100;
*/

/*
from WMCore/src/couchapps/WMStats/_attachments/js/minified/import-all-t0.min.js
    getWMBSTotalJobs: function() {
        return (this.getJobStatus("success") +
                this.getJobStatus("canceled") +
                this.getJobStatus( "transition") +
                this.getTotalFailure() +
                this.getTotalCooloff() +
                this.getTotalPaused() +
                this.getTotalQueued() +
                this.getRunning() +
                this.getPending());
    },

    getTotalFailure: function() {
        return (this.getJobStatus("failure.create") +
                this.getJobStatus("failure.submit") +
                this.getJobStatus("failure.exception"));
    },
*/

/*
from WMCore/src/couchapps/WMStats/_attachments/js/Views/Tables/T1/WMStats.SiteSummaryTable.js
              "title": "failure rate",
              "render": function (data, type, row, meta) {
                            var failJobs =  row.summary.getTotalFailure();
                            var successJobs = row.summary.getJobStatus("success");
                            var totalCompleteJobs = (successJobs + failJobs) || 1;
                            var result = failJobs / totalCompleteJobs * 100;
                            return (result.toFixed(1)  + "%");
                          }
            }
*/
