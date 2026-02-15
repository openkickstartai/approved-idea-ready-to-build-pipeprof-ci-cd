# PipeProf — CI/CD Pipeline Performance Analyzer

> Find every wasted minute and dollar in your CI/CD pipelines.

PipeProf analyzes your GitHub Actions / GitLab CI workflow files and identifies performance bottlenecks, cost waste, and optimization opportunities.

## 🚀 Quick Start

```bash
# Install
go install github.com/pipeprof/pipeprof@latest

# Or build from source
git clone https://github.com/pipeprof/pipeprof && cd pipeprof
go build -o pipeprof .

# Analyze a workflow
./pipeprof .github/workflows/ci.yml

# JSON output for CI integration
./pipeprof .github/workflows/ci.yml --json
```

## 📊 What It Detects

| Check | Severity | Typical Savings |
|---|---|---|
| Missing dependency caching | HIGH | ~30% build time |
| Expensive runner usage | HIGH | ~15% cost |
| No timeout-minutes set | WARN | Prevents runaway bills |
| No concurrency control | WARN | ~10% redundant runs |
| Unnecessary sequential jobs | WARN | ~20% wall time |

## 📊 Why Pay for PipeProf?

Free CLI catches the basics. **Pro** connects to your GitHub/GitLab API to analyze **real run durations**, track trends over time, and alert your team on Slack when costs spike. Teams using Pro save **$200–$2,000/month** on CI bills.

## 💰 Pricing

| Feature | Free (CLI) | Pro ($29/mo) | Enterprise ($99/mo) |
|---|---|---|---|
| YAML static analysis | ✅ | ✅ | ✅ |
| JSON export | ✅ | ✅ | ✅ |
| Workflows analyzed | 3/day | Unlimited | Unlimited |
| Live GitHub/GitLab API | ❌ | ✅ | ✅ |
| Historical cost tracking | ❌ | ✅ | ✅ |
| Duration-based analysis | ❌ | ✅ | ✅ |
| Team dashboard | ❌ | ❌ | ✅ |
| Slack/Teams alerts | ❌ | ❌ | ✅ |
| SSO & audit log | ❌ | ❌ | ✅ |
| Priority support | ❌ | ❌ | ✅ |

## License

MIT — Free CLI forever. Pro & Enterprise at [pipeprof.dev](https://pipeprof.dev).
