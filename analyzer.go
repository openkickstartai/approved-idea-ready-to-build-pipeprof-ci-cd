package main

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Name string         `yaml:"name"`
	Jobs map[string]Job `yaml:"jobs"`
}

type Job struct {
	RunsOn  string      `yaml:"runs-on"`
	Steps   []Step      `yaml:"steps"`
	Needs   interface{} `yaml:"needs"`
	Timeout int         `yaml:"timeout-minutes"`
}

type Step struct {
	Name string `yaml:"name"`
	Uses string `yaml:"uses"`
	Run  string `yaml:"run"`
}

type Issue struct {
	Severity string  `json:"severity"`
	Title    string  `json:"title"`
	Fix      string  `json:"fix"`
	Savings  float64 `json:"savings_per_month"`
}

type Report struct {
	WorkflowName     string  `json:"workflow_name"`
	TotalJobs        int     `json:"total_jobs"`
	Issues           []Issue `json:"issues"`
	MonthlyCost      float64 `json:"monthly_cost"`
	PotentialSavings float64 `json:"potential_savings"`
}

const (
	pricePerMin  = 0.008
	runsPerMonth = 200
	minsPerJob   = 5
)

func Analyze(data []byte) Report {
	var wf Workflow
	yaml.Unmarshal(data, &wf)
	r := Report{
		WorkflowName: wf.Name,
		TotalJobs:    len(wf.Jobs),
		Issues:       []Issue{},
	}
	r.MonthlyCost = float64(len(wf.Jobs)) * minsPerJob * pricePerMin * runsPerMonth
	hasCache := false
	raw := string(data)
	for name, job := range wf.Jobs {
		for _, step := range job.Steps {
			if strings.Contains(step.Uses, "cache") {
				hasCache = true
			}
		}
		if job.Timeout == 0 {
			r.addIssue("WARN", fmt.Sprintf("Job '%s' has no timeout-minutes", name),
				"Add timeout-minutes to prevent runaway builds", r.MonthlyCost*0.05)
		}
		if strings.Contains(job.RunsOn, "large") {
			r.addIssue("HIGH", fmt.Sprintf("Job '%s' uses expensive runner: %s", name, job.RunsOn),
				"Switch to standard runners or self-hosted", r.MonthlyCost*0.15)
		}
	}
	if !hasCache {
		r.addIssue("HIGH", "No dependency caching detected",
			"Add actions/cache to cut build time ~30%", r.MonthlyCost*0.3)
	}
	if !strings.Contains(raw, "concurrency") {
		r.addIssue("WARN", "No concurrency control",
			"Add concurrency group to cancel redundant runs", r.MonthlyCost*0.1)
	}
	depJobs := 0
	for _, job := range wf.Jobs {
		if job.Needs != nil {
			depJobs++
		}
	}
	if depJobs >= len(wf.Jobs)-1 && len(wf.Jobs) > 2 {
		r.addIssue("WARN", "Pipeline is mostly sequential",
			"Parallelize independent jobs to reduce wall time", r.MonthlyCost*0.2)
	}
	return r
}

func (r *Report) addIssue(sev, title, fix string, savings float64) {
	r.Issues = append(r.Issues, Issue{sev, title, fix, savings})
	r.PotentialSavings += savings
}
