package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof" // Register standard pprof endpoints
	"os"
	"sort"
	"strings"
	"sync"
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

// bufferPool reuses bytes.Buffer instances to eliminate heap churn during batch formatting.
var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

// EventProcessor handles high-throughput log aggregation with minimal allocations.
type EventProcessor struct {
	reportBuilder strings.Builder
}

func NewEventProcessor() *EventProcessor {
	return &EventProcessor{}
}

// ProcessBatch parses a slice of raw JSON log records and returns an aggregated SummaryReport.
// Memory allocations are minimized via slice preallocation, buffer pooling, and avoiding string concatenation.
func (p *EventProcessor) ProcessBatch(records [][]byte) (SummaryReport, error) {
	if len(records) == 0 {
		return SummaryReport{}, ErrEmptyBatch
	}

	// Preallocate maps and structures based on batch size hint
	serviceMap := make(map[string]*ServiceStat, 16)
	totalErrors := 0
	validEvents := 0

	// Acquire pooled buffer for scratch JSON decoding / verification
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

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

	// Preallocate slice with exact capacity to eliminate reallocation copying
	services := make([]ServiceStat, 0, len(serviceMap))
	for _, stat := range serviceMap {
		services = append(services, *stat)
	}

	// Deterministic ordering by service name
	sort.Slice(services, func(i, j int) bool {
		return services[i].Service < services[j].Service
	})

	// Format compact digest using strings.Builder with pre-allocated growth
	var b strings.Builder
	b.Grow(64 + len(services)*32)
	b.WriteString("TOTAL=")
	b.WriteString(fmt.Sprintf("%d", validEvents))
	b.WriteString(";ERRORS=")
	b.WriteString(fmt.Sprintf("%d", totalErrors))
	b.WriteString(";SERVICES=")

	for i, s := range services {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(s.Service)
		b.WriteByte(':')
		b.WriteString(fmt.Sprintf("%d", s.Total))
	}

	report := SummaryReport{
		TotalEvents:   validEvents,
		TotalErrors:   totalErrors,
		Services:      services,
		GeneratedAt:   time.Now().UTC(),
		DigestSummary: b.String(),
	}

	return report, nil
}

// FormatReportJSON serializes the SummaryReport using a pooled buffer directly to the writer.
func FormatReportJSON(w io.Writer, report SummaryReport) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// StartPprofServer starts an HTTP server on the given address hosting pprof profiling endpoints.
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
	pprofAddr := ":6060"
	fmt.Printf("Starting pprof profiling endpoint on %s\n", pprofAddr)
	srv := StartPprofServer(pprofAddr)
	defer srv.Close()

	processor := NewEventProcessor()

	// Demonstration batch
	sampleRecords := [][]byte{
		[]byte(`{"timestamp":"2026-10-23T12:00:00Z","level":"INFO","service":"auth","message":"token issued"}`),
		[]byte(`{"timestamp":"2026-10-23T12:00:01Z","level":"ERROR","service":"auth","message":"invalid credentials"}`),
		[]byte(`{"timestamp":"2026-10-23T12:00:02Z","level":"INFO","service":"billing","message":"payment processed"}`),
	}

	report, err := processor.ProcessBatch(sampleRecords)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing batch: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processed batch successfully:")
	_ = FormatReportJSON(os.Stdout, report)
}
