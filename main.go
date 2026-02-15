package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("PipeProf — CI/CD Pipeline Performance Analyzer")
		fmt.Println("Usage: pipeprof <workflow.yml> [--json]")
		fmt.Println("       pipeprof analyze <workflow.yml> [--json]")
		os.Exit(0)
	}
	jsonMode := false
	path := ""
	for _, a := range args {
		switch {
		case a == "--json":
			jsonMode = true
		case a != "analyze":
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
	if jsonMode {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(report)
	} else {
		printReport(report)
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
