package search

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	vcpuRegex   = regexp.MustCompile(`(?i)\b(\d+)\s*(?:vcpu|cpu|core)s?\b`)
	ramRegex    = regexp.MustCompile(`(?i)\b(\d+)\s*(?:gb|gib|g)\b`)
	priceRegex  = regexp.MustCompile(`(?i)(?:under|<|<=|\$)\s*(\d+(?:\.\d+)?)\s*(?:usd|\$)?\b`)
	nvmeRegex   = regexp.MustCompile(`(?i)\bnvme\b`)
	ssdRegex    = regexp.MustCompile(`(?i)\bssd\b`)
	hddRegex    = regexp.MustCompile(`(?i)\bhdd\b`)
	gpuRegex    = regexp.MustCompile(`(?i)\b(gpu|rtx|a100|h100|vram)\b`)
	bareMetalRe = regexp.MustCompile(`(?i)\b(bare\s*metal|dedicated\s*server)\b`)
	dediCPURe   = regexp.MustCompile(`(?i)\b(dedicated\s*cpu|cpu-optimized)\b`)
)

var knownRegions = map[string]string{
	"frankfurt": "fra",
	"germany":   "fra",
	"tokyo":     "nrt",
	"japan":     "nrt",
	"seoul":     "icn",
	"korea":     "icn",
	"new york":  "nyc",
	"nyc":       "nyc",
	"london":    "lon",
	"uk":        "lon",
	"singapore": "sgp",
	"helsinki":  "hel",
	"finland":   "hel",
}

// ParseNLQuery extracts structured parameters from natural language input.
func ParseNLQuery(q string) (ParsedQueryMetadata, SearchRequest) {
	var meta ParsedQueryMetadata
	var req SearchRequest
	lower := strings.ToLower(q)

	// 1. vCPU extraction
	if matches := vcpuRegex.FindStringSubmatch(lower); len(matches) > 1 {
		if v, err := strconv.Atoi(matches[1]); err == nil && v > 0 {
			meta.DetectedVCPU = v
			req.MinVCPU = v
		}
	}

	// 2. RAM extraction
	if matches := ramRegex.FindStringSubmatch(lower); len(matches) > 1 {
		if v, err := strconv.Atoi(matches[1]); err == nil && v > 0 {
			meta.DetectedRAMGB = v
			req.MinRAMMB = int64(v) * 1024
		}
	}

	// 3. Storage Type
	if nvmeRegex.MatchString(lower) {
		meta.DetectedDiskType = "NVME"
		req.DiskType = "NVME"
	} else if ssdRegex.MatchString(lower) {
		meta.DetectedDiskType = "SSD"
		req.DiskType = "SSD"
	} else if hddRegex.MatchString(lower) {
		meta.DetectedDiskType = "HDD"
		req.DiskType = "HDD"
	}

	// 4. Max Price extraction (e.g. under $40, < 40)
	if strings.Contains(lower, "under") || strings.Contains(lower, "<") {
		if matches := priceRegex.FindStringSubmatch(lower); len(matches) > 1 {
			if d, err := decimal.NewFromString(matches[1]); err == nil && d.GreaterThan(decimal.Zero) {
				meta.DetectedMaxPrice = d
				req.MaxMonthlyPriceUSD = d
			}
		}
	}

	// 5. Category detection
	if gpuRegex.MatchString(lower) {
		meta.DetectedCategory = "GPU_SERVER"
		req.Category = "GPU_SERVER"
	} else if bareMetalRe.MatchString(lower) {
		meta.DetectedCategory = "BARE_METAL"
		req.Category = "BARE_METAL"
	} else if dediCPURe.MatchString(lower) {
		meta.DetectedCategory = "DEDICATED_CPU_VM"
		req.Category = "DEDICATED_CPU_VM"
	}

	// 6. Region mapping
	for name, code := range knownRegions {
		if strings.Contains(lower, name) {
			meta.DetectedRegion = code
			req.Region = code
			break
		}
	}

	req.Query = q
	return meta, req
}

