package main

import "testing"

func TestAnalyzeSTAR_FullResponse(t *testing.T) {
	response := `When I was working on the payment team, we had a critical issue with timeouts. I needed to fix it before the launch deadline. I implemented a circuit breaker pattern using Go channels. This reduced latency by 60%.`
	result := AnalyzeSTAR(response)
	if result.Score < 3 {
		t.Errorf("expected at least 3 STAR elements, got %d", result.Score)
	}
}

func TestAnalyzeSTAR_MinimalResponse(t *testing.T) {
	response := `I built a thing.`
	result := AnalyzeSTAR(response)
	if result.Score > 1 {
		t.Errorf("expected 0-1 STAR elements for minimal response, got %d", result.Score)
	}
}

func TestAnalyzeSTAR_AllElements(t *testing.T) {
	response := `When I was working at the company (situation), I needed to solve a scaling problem (task). I designed a new architecture (action). This improved performance by 50% (result).`
	result := AnalyzeSTAR(response)
	if result.Score != 4 {
		t.Errorf("expected 4 STAR elements, got %d", result.Score)
	}
}

func TestAnalyzeSTAR_ActionKeywords(t *testing.T) {
	response := `I implemented a new caching layer. I designed the schema. I built the API.`
	result := AnalyzeSTAR(response)
	actionFound := false
	for _, c := range result.Checks {
		if c.Element == Action && c.Found {
			actionFound = true
		}
	}
	if !actionFound {
		t.Errorf("expected Action element to be found")
	}
}

func TestAnalyzeSTAR_ResultKeywords(t *testing.T) {
	response := `The outcome was that we reduced costs by 30% and improved reliability.`
	result := AnalyzeSTAR(response)
	resultFound := false
	for _, c := range result.Checks {
		if c.Element == Result && c.Found {
			resultFound = true
		}
	}
	if !resultFound {
		t.Errorf("expected Result element to be found")
	}
}

func TestAnalyzeSTAR_AllElementsPresentInOutput(t *testing.T) {
	result := AnalyzeSTAR("test")
	if len(result.Checks) != 4 {
		t.Errorf("expected 4 checks, got %d", len(result.Checks))
	}
}
