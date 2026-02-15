package main

import (
	"strings"
	"testing"
)

const wfNoCache = `name: Build
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm install
      - run: npm test
  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - run: echo deploy
`

const wfOptimized = `name: Optimized
on: push
concurrency:
  group: ci-${{ github.ref }}
jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@v4
      - uses: actions/cache@v3
      - run: npm test
`

const wfExpensive = `name: Heavy
on: push
jobs:
  build:
    runs-on: ubuntu-latest-4xlarge
    steps:
      - uses: actions/checkout@v4
      - run: make build
  test:
    needs: build
    runs-on: ubuntu-latest-xlarge
    steps:
      - run: make test
  lint:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - run: make lint
`

func TestDetectsMissingCache(t *testing.T) {
	r := Analyze([]byte(wfNoCache))
	if r.TotalJobs != 2 {
		t.Errorf("expected 2 jobs, got %d", r.TotalJobs)
	}
	found := false
	for _, iss := range r.Issues {
		if iss.Title == "No dependency caching detected" {
			found = true
		}
	}
	if !found {
		t.Error("should detect missing cache")
	}
	if r.PotentialSavings <= 0 {
		t.Error("expected positive potential savings")
	}
}

func TestOptimizedHasFewerIssues(t *testing.T) {
	r := Analyze([]byte(wfOptimized))
	if r.TotalJobs != 1 {
		t.Errorf("expected 1 job, got %d", r.TotalJobs)
	}
	for _, iss := range r.Issues {
		if iss.Title == "No dependency caching detected" {
			t.Error("should not flag cache when actions/cache is used")
		}
		if iss.Title == "No concurrency control" {
			t.Error("should not flag concurrency when present")
		}
	}
	if r.MonthlyCost <= 0 {
		t.Error("expected positive monthly cost")
	}
}

func TestDetectsExpensiveRunners(t *testing.T) {
	r := Analyze([]byte(wfExpensive))
	highCount := 0
	for _, iss := range r.Issues {
		if iss.Severity == "HIGH" && strings.Contains(iss.Title, "expensive runner") {
			highCount++
		}
	}
	if highCount < 2 {
		t.Errorf("expected at least 2 expensive runner issues, got %d", highCount)
	}
	if r.TotalJobs != 3 {
		t.Errorf("expected 3 jobs, got %d", r.TotalJobs)
	}
	seqFound := false
	for _, iss := range r.Issues {
		if iss.Title == "Pipeline is mostly sequential" {
			seqFound = true
		}
	}
	if !seqFound {
		t.Error("should detect sequential pipeline with 3 chained jobs")
	}
}
