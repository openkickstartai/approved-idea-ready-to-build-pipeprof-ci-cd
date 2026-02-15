package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("PipeProf — CI/CD Pipeline Performance Analyzer")
		fmt.Println("Usage: pipeprof <workflow.yml> [--format table|json]")
		fmt.Println("       pipeprof analyze <workflow.yml> [--format table|json]")
		os.Exit(0)
	}
	format := "table"
	path := ""
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--json":
			format = "json"
		case a == "--format" && i+1 < len(args):
			format = args[i+1]
			i++
		case strings.HasPrefix(a, "--format="):
			format = a[len("--format="):]
		case a != "analyze" && !strings.HasPrefix(a, "--"):
			path = a
		}
	}
	if path == "" {
		fmt.Fprintln(os.Stderr, "Error: no workflow file specified")
		os.Exit(1)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	report := Analyze(data)
	analysisReport := BuildAnalysisReport(data, report)

	switch format {
	case "json":
		FormatJSON(analysisReport, os.Stdout)
	case "table":
		FormatTable(analysisReport, os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown format %q (use 'table' or 'json')\n", format)
		os.Exit(1)
	}
}

func printReport(r Report) {
	fmt.Println("══════════════════════════════════════════════")
	fmt.Printf("  PipeProf Report: %s\n", r.WorkflowName)
	fmt.Println("══════════════════════════════════════════════")
	fmt.Printf("  Jobs: %d  |  Est. cost: $%.2f/mo\n", r.TotalJobs, r.MonthlyCost)
	fmt.Printf("  Issues: %d  |  Savings: $%.2f/mo\n", len(r.Issues), r.PotentialSavings)
	fmt.Println("──────────────────────────────────────────────")
	for i, iss := range r.Issues {
		fmt.Printf("  %d. [%s] %s\n", i+1, iss.Severity, iss.Title)
		fmt.Printf("     → %s (saves $%.2f/mo)\n\n", iss.Fix, iss.Savings)
	}
	if len(r.Issues) == 0 {
		fmt.Println("  ✅ Pipeline looks optimized!")
	}
	fmt.Println("──────────────────────────────────────────────")
	fmt.Println("  🚀 Upgrade to Pro: https://pipeprof.dev/pro")
	fmt.Println("══════════════════════════════════════════════")
}
