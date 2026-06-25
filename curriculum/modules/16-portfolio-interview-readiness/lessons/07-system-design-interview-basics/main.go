package main

import (
	"fmt"
	"math"
	"strings"
)

type DataUnit int

const (
	Bytes DataUnit = iota
	Kilobytes
	Megabytes
	Gigabytes
	Terabytes
	RequestsPerSecond
	Records
)

func (d DataUnit) String() string {
	switch d {
	case Bytes:
		return "bytes"
	case Kilobytes:
		return "KB"
	case Megabytes:
		return "MB"
	case Gigabytes:
		return "GB"
	case Terabytes:
		return "TB"
	case RequestsPerSecond:
		return "req/s"
	case Records:
		return "records"
	default:
		return "units"
	}
}

type Requirement struct {
	Name        string
	Value       float64
	Unit        DataUnit
	Description string
}

type CapacityEstimate struct {
	Requirement     Requirement
	EstimatedValue  float64
	EstimatedUnit   DataUnit
	Assumptions     []string
	Recommendations []string
}

type SystemDesign struct {
	Name         string
	Requirements []Requirement
	Estimates    []CapacityEstimate
}

func (sd *SystemDesign) AddRequirement(name string, value float64, unit DataUnit, desc string) {
	sd.Requirements = append(sd.Requirements, Requirement{name, value, unit, desc})
}

func bytesToUnit(bytes float64) (float64, DataUnit) {
	switch {
	case bytes >= 1<<40:
		return bytes / (1 << 40), Terabytes
	case bytes >= 1<<30:
		return bytes / (1 << 30), Gigabytes
	case bytes >= 1<<20:
		return bytes / (1 << 20), Megabytes
	case bytes >= 1<<10:
		return bytes / (1 << 10), Kilobytes
	default:
		return bytes, Bytes
	}
}

func (sd *SystemDesign) EstimateStorage() {
	var totalBytes float64
	for _, req := range sd.Requirements {
		if req.Unit == Records || req.Unit == RequestsPerSecond {
			continue
		}
		switch req.Unit {
		case Bytes:
			totalBytes += req.Value
		case Kilobytes:
			totalBytes += req.Value * (1 << 10)
		case Megabytes:
			totalBytes += req.Value * (1 << 20)
		case Gigabytes:
			totalBytes += req.Value * (1 << 30)
		case Terabytes:
			totalBytes += req.Value * (1 << 40)
		}
	}
	estimated, unit := bytesToUnit(totalBytes)
	sd.Estimates = append(sd.Estimates, CapacityEstimate{
		Requirement:     Requirement{Name: "Total Storage", Unit: Bytes},
		EstimatedValue:  estimated,
		EstimatedUnit:   unit,
		Assumptions:     []string{"assumes single-copy storage", "no compression applied"},
		Recommendations: []string{"add replication factor (3x for durability)", "consider compression (2-5x reduction)", "plan for 20% annual growth"},
	})
}

func (sd *SystemDesign) EstimateThroughput() {
	var totalRPS float64
	growth := 1.20
	for _, req := range sd.Requirements {
		if req.Unit == RequestsPerSecond {
			totalRPS += req.Value
		}
	}
	peakRPS := totalRPS * 2.0

	sd.Estimates = append(sd.Estimates, CapacityEstimate{
		Requirement:     Requirement{Name: "Throughput", Value: totalRPS, Unit: RequestsPerSecond},
		EstimatedValue:  peakRPS,
		EstimatedUnit:   RequestsPerSecond,
		Assumptions:     []string{"peak traffic is 2x average", fmt.Sprintf("annual growth rate: %.0f%%", (growth-1)*100)},
		Recommendations: []string{"design for 2x peak headroom", "use connection pooling", "consider CDN for static content"},
	})

	storageRPS := totalRPS * 365 * 86400
	for _, req := range sd.Requirements {
		if req.Unit == Records {
			storageRPS *= req.Value
		}
	}
	_ = storageRPS
}

func (sd *SystemDesign) PrintReport() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("System Design: %s\n", sd.Name))
	b.WriteString(strings.Repeat("-", 60) + "\n")

	b.WriteString("Requirements:\n")
	for _, req := range sd.Requirements {
		val := req.Value
		if val == math.Trunc(val) {
			b.WriteString(fmt.Sprintf("  - %s: %.0f %s\n", req.Name, val, req.Unit))
		} else {
			b.WriteString(fmt.Sprintf("  - %s: %.2f %s\n", req.Name, val, req.Unit))
		}
		if req.Description != "" {
			b.WriteString(fmt.Sprintf("    %s\n", req.Description))
		}
	}

	b.WriteString("\nCapacity Estimates:\n")
	for _, est := range sd.Estimates {
		b.WriteString(fmt.Sprintf("  - %s: %.0f %s\n", est.Requirement.Name, est.EstimatedValue, est.EstimatedUnit))
		b.WriteString("    Assumptions:\n")
		for _, a := range est.Assumptions {
			b.WriteString(fmt.Sprintf("      * %s\n", a))
		}
		b.WriteString("    Recommendations:\n")
		for _, r := range est.Recommendations {
			b.WriteString(fmt.Sprintf("      * %s\n", r))
		}
	}

	return b.String()
}

func main() {
	design := SystemDesign{Name: "URL Shortener"}

	design.AddRequirement("Daily Active Users", 1000000, RequestsPerSecond, "expected MAU after first year")
	design.AddRequirement("Write RPS", 115, RequestsPerSecond, "short URL creations")
	design.AddRequirement("Read RPS", 11500, RequestsPerSecond, "redirect requests (100:1 read/write ratio)")
	design.AddRequirement("URL Record Size", 512, Bytes, "average entry including metadata")
	design.AddRequirement("Total URL Records", 36500000*5, Records, "36.5M URLs/year x 5 years")

	design.EstimateStorage()
	design.EstimateThroughput()

	fmt.Print(design.PrintReport())
}
