package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"
)

var (
	ErrEmptyBatch  = errors.New("empty batch provided")
	ErrCorruptData = errors.New("corrupt log record")
)

type LogEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Service   string    `json:"service"`
	Message   string    `json:"message"`
}

type ServiceStat struct {
	Service string `json:"service"`
	Total   int    `json:"total"`
	Errors  int    `json:"errors"`
}

type SummaryReport struct {
	TotalEvents   int           `json:"total_events"`
	TotalErrors   int           `json:"total_errors"`
	Services      []ServiceStat `json:"services"`
	GeneratedAt   time.Time     `json:"generated_at"`
	DigestSummary string        `json:"digest_summary"`
}

// TODO: Task 1 - Eliminate Heap Churn with sync.Pool
// Define a sync.Pool for bytes.Buffer to reuse scratch buffers during JSON encoding.

type EventProcessor struct {
	// TODO: Task 2 - Pre-size string builder or cache reusable structures
}

func NewEventProcessor() *EventProcessor {
	return &EventProcessor{}
}

// ProcessBatch parses a slice of raw JSON log records and returns an aggregated SummaryReport.
// UNOPTIMIZED: Incurs heavy GC pressure due to unpooled allocations and dynamic slice growth.
func (p *EventProcessor) ProcessBatch(records [][]byte) (SummaryReport, error) {
	if len(records) == 0 {
		return SummaryReport{}, ErrEmptyBatch
	}

	// UNOPTIMIZED: Map with no initial capacity hint
	serviceMap := make(map[string]*ServiceStat)
	totalErrors := 0
	validEvents := 0

	for _, rec := range records {
		if len(rec) == 0 {
			continue
		}

		var ev LogEvent
		if err := json.Unmarshal(rec, &ev); err != nil {
			return SummaryReport{}, fmt.Errorf("%w: %v", ErrCorruptData, err)
		}

		validEvents++
		isErr := ev.Level == "ERROR" || ev.Level == "FATAL"
		if isErr {
			totalErrors++
		}

		stat, exists := serviceMap[ev.Service]
		if !exists {
			stat = &ServiceStat{
				Service: ev.Service,
				Total:   0,
				Errors:  0,
			}
			serviceMap[ev.Service] = stat
		}

		stat.Total++
		if isErr {
			stat.Errors++
		}
	}

	// UNOPTIMIZED: Slice append with zero initial capacity triggers multiple re-allocations
	var services []ServiceStat
	for _, stat := range serviceMap {
		services = append(services, *stat)
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].Service < services[j].Service
	})

	// UNOPTIMIZED: Naive string concatenation triggers O(N^2) memory copying on the heap
	// TODO: Task 2 - Replace string concatenation with strings.Builder and pre-allocated Grow()
	digest := "TOTAL=" + fmt.Sprintf("%d", validEvents) + ";ERRORS=" + fmt.Sprintf("%d", totalErrors) + ";SERVICES="
	for i, s := range services {
		if i > 0 {
			digest += ","
		}
		digest += s.Service + ":" + fmt.Sprintf("%d", s.Total)
	}

	report := SummaryReport{
		TotalEvents:   validEvents,
		TotalErrors:   totalErrors,
		Services:      services,
		GeneratedAt:   time.Now().UTC(),
		DigestSummary: digest,
	}

	return report, nil
}

// FormatReportJSON serializes the SummaryReport.
// UNOPTIMIZED: Allocates fresh buffers on every call.
func FormatReportJSON(w io.Writer, report SummaryReport) error {
	// TODO: Task 1 - Use sync.Pool for buffer reuse
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}

// TODO: Task 3 - Mount pprof profiling endpoints
// Import _ "net/http/pprof" and expose HTTP server with /debug/pprof/
func StartPprofServer(addr string) *http.Server {
	srv := &http.Server{
		Addr: addr,
	}
	go func() {
		_ = srv.ListenAndServe()
	}()
	return srv
}

func main() {
	processor := NewEventProcessor()

	sampleRecords := [][]byte{
		[]byte(`{"timestamp":"2026-10-23T12:00:00Z","level":"INFO","service":"auth","message":"token issued"}`),
		[]byte(`{"timestamp":"2026-10-23T12:00:01Z","level":"ERROR","service":"auth","message":"invalid credentials"}`),
	}

	report, err := processor.ProcessBatch(sampleRecords)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	_ = FormatReportJSON(os.Stdout, report)
}
