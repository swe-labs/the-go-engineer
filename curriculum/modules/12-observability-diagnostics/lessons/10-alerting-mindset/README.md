# Alerting mindset

## Learning objective

Design meaningful SLO-based alerts by calculating SLIs, burn rates, and alert severity, and distinguish actionable alerts from noise.

## Why this matters

Bad alerts destroy on-call trust. A team that receives 50 alerts per night ignores all of them — including the one that signals a real outage. Alert fatigue is the leading cause of incident response delays. Google SRE research shows that teams with >10 alerts per on-call shift have a mean-time-to-respond (MTTR) 3x higher than teams with <5 alerts. Mastering alerting mindset means designing alerts that fire only when a human must act. Everything else is a dashboard or a log.

## Mental model

Alerting is a smoke detector, not a thermostat. A thermostat constantly adjusts (automated remediation). A smoke detector only makes noise when there is a fire (user-impacting problem). If the smoke detector goes off every time someone cooks (deployment causes a brief latency spike), the occupants disable it (alert fatigue). The correct design: the smoke detector only goes off for actual fires — situations that require immediate human action. Everything else is a dashboard or a log.

The key metric is the _burn rate_: how fast are you consuming your error budget? If your SLO is 99.9% availability (monthly error budget: 43 minutes of downtime), and you've consumed 10 minutes in the last hour, your burn rate is (10 min / 1 hour) / (43 min / 30 days) = 168x — you are burning budget 168 times faster than allowed. A burn rate of > 10x for 1 hour or > 2x for 6 hours triggers a page.

The analogy breaks because smoke detectors are binary (fire/no fire), while alerts have severity levels and can trigger automated responses before paging a human.

## Core idea

Three fundamental concepts underpin every alert:

- **SLI (Service Level Indicator)**: a quantitative measure of a service property. Examples: request latency (p99 < 200ms), error rate (< 1%), throughput (> 1000 req/s).
- **SLO (Service Level Objective)**: a target value or range for an SLI over a time window. Example: "p99 latency < 200ms for 99.9% of requests over a 30-day window."
- **SLA (Service Level Agreement)**: a contractual commitment to meet an SLO, with consequences (credits, penalties) for violation.

Alert severity levels:

| Severity | Meaning | Response time | Example |
|---|---|---|---|
| P0/Critical | Service unavailable or data loss | Immediate (<5 min) | Payment service down |
| P1/High | Major feature degradation | <30 min | Search latency > 5s |
| P2/Medium | Partial degradation, no user impact | <4 hours | Dashboard slow to load |
| P3/Low | Cosmetic, non-urgent | Next business day | Minor UI glitch |

## Under the hood

A burn-rate alert evaluates PromQL queries at the evaluation interval. For a 99.9% SLO with a 1-hour window:

```
# Alert if burn rate exceeds 10x in 1 hour
- alert: HighErrorRate
  expr: |
    (
      1 - (sum(rate(http_requests_total{status=~"5.."}[1h]))
           / sum(rate(http_requests_total[1h])))
    ) < 0.999
    and
    (
      sum(rate(http_requests_total[1h])) > 100
    )
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "Error rate SLO breach (burn rate {{ $value }})"
```

The PromQL engine computes rate queries over the specified window. The `for: 5m` clause ensures the condition is sustained for 5 minutes before firing — preventing flapping from transient spikes. Alertmanager receives the firing alert, applies grouping (e.g., group by service), silencing (maintenance windows), and inhibition (suppress low-severity when critical fires), then sends the notification via the configured receiver (PagerDuty, Slack, email).

An _alert_ in Alertmanager goes through states: Inactive → Pending (breach detected but `for` not elapsed) → Firing (breach sustained for `for` duration) → Resolved (breach cleared).

## How Go uses it

Alerting logic is typically expressed in Prometheus alerting rules (YAML configuration) rather than Go code. However, Go services often need to:

1. Expose SLI metrics via Prometheus (latency histogram, error counter).
2. Compute burn rates in code for dynamic alert thresholds.
3. Build alert simulation tools to test rule effectiveness before deploying.
4. Implement health check endpoints that include SLI data for load balancer decision-making.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

type SLOConfig struct {
	Window        time.Duration
	Target        float64 // e.g., 0.999 for 99.9%
	ErrorBudget   float64 // computed: 1 - target
}

type BurnRateAlert struct {
	Severity     string
	Threshold    float64 // burn rate multiplier
	Window       time.Duration
}

type SLI struct {
	TotalRequests int
	ErrorRequests int
}

func (s SLI) ErrorRate() float64 {
	if s.TotalRequests == 0 {
		return 0
	}
	return float64(s.ErrorRequests) / float64(s.TotalRequests)
}

func (s SLI) SuccessRate() float64 {
	return 1 - s.ErrorRate()
}

func computeBurnRate(currentErrors int, currentTotal int, since time.Duration, slo SLOConfig) float64 {
	actualErrorRate := float64(currentErrors) / float64(currentTotal)
	maxAllowedErrorRate := 1 - slo.Target

	if maxAllowedErrorRate == 0 {
		return 0
	}
	return actualErrorRate / maxAllowedErrorRate
}

func evaluateAlert(sli SLI, since time.Duration, slo SLOConfig, alerts []BurnRateAlert) []string {
	var firing []string
	burnRate := computeBurnRate(sli.ErrorRequests, sli.TotalRequests, since, slo)

	for _, alert := range alerts {
		if burnRate >= alert.Threshold {
			firing = append(firing, alert.Severity)
		}
	}
	return firing
}

func main() {
	slo := SLOConfig{
		Window:      30 * 24 * time.Hour,
		Target:      0.999,
		ErrorBudget: 0.001,
	}

	alerts := []BurnRateAlert{
		{Severity: "critical", Threshold: 10, Window: 1 * time.Hour},
		{Severity: "warning", Threshold: 2, Window: 6 * time.Hour},
	}

	scenarios := []struct {
		name         string
		total        int
		errors       int
		since        time.Duration
	}{
		{"normal operation", 10000, 2, 1 * time.Hour},
		{"brief spike", 1000, 50, 1 * time.Hour},
		{"sustained degradation", 100000, 500, 6 * time.Hour},
		{"major outage", 5000, 200, 30 * time.Minute},
	}

	fmt.Printf("SLO: %.1f%%, Error Budget: %.1f%%\n\n", slo.Target*100, slo.ErrorBudget*100)

	for _, s := range scenarios {
		sli := SLI{TotalRequests: s.total, ErrorRequests: s.errors}
		rate := sli.ErrorRate() * 100
		burn := computeBurnRate(s.errors, s.total, s.since, slo)
		firing := evaluateAlert(sli, s.since, slo, alerts)

		fmt.Printf("[%s]\n", s.name)
		fmt.Printf("  Requests: %d, Errors: %d, Error rate: %.2f%%\n", s.total, s.errors, rate)
		fmt.Printf("  Burn rate: %.1fx\n", burn)
		if len(firing) == 0 {
			fmt.Println("  Status: OK")
		} else {
			fmt.Printf("  Alerts firing: %v\n", firing)
		}
		fmt.Println()
	}
}
```

## Step-by-step execution

1. Define SLO with target (99.9%) and compute error budget (0.1%).
2. Define two burn-rate alert thresholds: critical at 10x burn, warning at 2x burn.
3. For each scenario, compute the error rate and burn rate.
4. Burn rate = actual error rate / max allowed error rate.
5. If burn rate >= threshold, the alert fires.
6. Normal operation (0.02% errors) has burn rate 0.2x — no alert.
7. Brief spike (5% errors) has burn rate 50x — critical alert fires.
8. Sustained degradation (0.5% errors over 6h) has burn rate 5x — warning alert fires.
9. Major outage (4% errors) has burn rate 40x — critical alert fires.

## Common mistakes

- Mistake: Alerting on every metric that can go wrong (CPU > 80%, memory > 90%, error count > 0).
  - Why it happens: Engineers treat monitoring and alerting as the same thing. They create alerts for every dashboard panel.
  - Fix: An alert must represent a user-impacting condition requiring human action within a specific time window. If the team ignores an alert, delete it or downgrade to a dashboard annotation.

- Mistake: Setting static thresholds without understanding the service's normal distribution.
  - Why it happens: Picking round numbers (500ms, 80% CPU) that have no relationship to the service's actual performance.
  - Fix: Measure baseline metrics for 2 weeks, then set thresholds at 3 standard deviations above the mean. Use dynamic thresholds where possible.

- Mistake: Creating alerts without a runbook.
  - Why it happens: The alert is created during an incident, and the runbook is "documented later" — which never happens.
  - Fix: Every alert must include a runbook URL in its annotations. The runbook must specify: (1) what to check, (2) how to diagnose, (3) how to mitigate, (4) who to escalate to.

## Debugging walkthrough

A team receives 3-5 PagerDuty alerts per night for "High CPU" on their Go service. The on-call engineer checks the dashboard, sees CPU at 85%, and acknowledges — there is nothing to do. The alert fires again the next night.

**Symptom**: Chronic "High CPU" alerts that never lead to any action.

**Investigation**:
1. Check the alert threshold: CPU > 80% for 5 minutes.
2. Check the service's normal CPU: baseline is 70-85%, spiking to 90% during peak traffic.
3. The threshold is set below the normal operating range — it fires daily.

**Root cause**: The alert threshold was set during service initialization before baseline data existed. The static threshold (80%) is below the service's actual operating range.

**Fix**: Change the threshold to 95% (the actual danger zone) or use dynamic thresholding. Better: replace the CPU alert with a burn-rate alert on the SLO. CPU > 95% is a dashboard concern; SLO breach is an actionable alert.

## Production notes

- Google SRE recommends <10 critical alerts per service. Each alert should fire <1x per quarter.
- Use multi-window burn-rate alerts: a fast window (1 hour at 10x burn) catches sudden outages, and a slow window (6 hours at 2x burn) catches gradual degradation.
- Always include a link to the runbook, dashboard, and recent change log in every alert notification.
- Test alerts regularly using fault injection (chaos engineering) to verify they fire correctly and the runbook is accurate.
- Use alert silences for planned maintenance. Never disable alerts permanently — fix the threshold or delete the alert.

## Performance implications

- Alert evaluation (PromQL queries) is CPU-intensive on the Prometheus server. A single complex `histogram_quantile` query can take seconds on a large TSDB. Use recording rules to pre-compute expensive queries.
- Burn-rate alerts are more computationally expensive than static threshold alerts but produce far fewer alerts, reducing the human cost of alert fatigue.
- In Go services, instrumenting metrics for SLI tracking adds negligible overhead (~10ns per counter increment). The cost of alerting is in the monitoring infrastructure, not the application.

## Practice task

Implement a `BudgetTracker` that:
1. Takes an SLO target (e.g., 99.9%) and a window (30 days) in the constructor.
2. Has a `Record(success bool)` method that records a successful or failed request.
3. Has a `BurnRate(window time.Duration) float64` method that computes the burn rate over the given window.
4. Has a `RemainingBudget() float64` method that returns the remaining error budget as a percentage.
5. Has a `ShouldAlert(burnRateThreshold float64) bool` method that returns true when the burn rate exceeds the threshold.

Write tests that verify budget consumption over simulated request streams.

## Tests / verification

```bash
go test ./curriculum/modules/12-observability-diagnostics/lessons/10-alerting-mindset -v
go run ./curriculum/modules/12-observability-diagnostics/lessons/10-alerting-mindset
```

## Review questions

1. What is the difference between an SLI, SLO, and SLA?
2. Why does a burn rate of 10x over 1 hour trigger a critical alert even if the monthly SLO is still compliant?
3. A team's service has 99.9% SLO. Over 24 hours, the error rate is 0.5%. How much error budget was consumed?
4. What is alert fatigue and how does burn-rate alerting reduce it?
5. Why should every alert include a runbook URL in its annotations?

## NEXT UP

Incident debugging — applying structured methodology to diagnose production incidents using observability data.
