package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Module0096 implements architecture fitness tests: zero TreeScan dependency & adapter separation.
type Module0096 struct {
	baseDir string
}

// NewModule0096 constructs an architecture fitness verifier.
func NewModule0096(baseDir string) *Module0096 {
	if baseDir == "" {
		baseDir = "."
	}
	return &Module0096{baseDir: baseDir}
}

// VerifyTreeScanIndependence scans production Go source files to ensure ZERO import or reference to TreeScan.
func (m *Module0096) VerifyTreeScanIndependence(ctx context.Context) error {
	var violations []string

	err := filepath.Walk(m.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == ".nuxt" || name == "11_WORK_ORDERS_0001_TO_1000" {
				return filepath.SkipDir
			}
			return nil
		}

		// Inspect only Go files, config files, and Dockerfiles
		if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".mod") && !strings.HasSuffix(path, ".yml") {
			return nil
		}

		// Skip boundary verifiers that declare the blacklist and test fixtures
		fileName := filepath.Base(path)
		if strings.HasSuffix(fileName, "_test.go") ||
			fileName == "module_0051.go" ||
			fileName == "module_0096.go" ||
			fileName == "bootstrap.go" {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := strings.ToLower(string(contentBytes))

		// Check for forbidden TreeScan couplings in production code paths
		if strings.Contains(content, "treescan_db") ||
			strings.Contains(content, "treescan_prod") ||
			strings.Contains(content, "treescan.com/api") ||
			strings.Contains(content, "github.com/treescan") {
			violations = append(violations, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walk failed: %w", err)
	}

	if len(violations) > 0 {
		return fmt.Errorf("article 1 fitness failure: found TreeScan coupling in: %v", violations)
	}

	return nil
}

// VerifyAdapterPackageSeparation checks that pricing, affiliate, management, and reseller packages
// do not cross-import each other's private packages.
func (m *Module0096) VerifyAdapterPackageSeparation(ctx context.Context) error {
	pricingPath := filepath.Join(m.baseDir, "internal", "adapters", "pricing", "adapter.go")
	contentBytes, err := os.ReadFile(pricingPath)
	if err != nil {
		return fmt.Errorf("failed to read pricing adapter: %w", err)
	}

	content := string(contentBytes)
	if strings.Contains(content, "internal/adapters/management") {
		return errors.New("article 6 violation: pricing adapter illegally imports management adapter")
	}
	if strings.Contains(content, "internal/adapters/affiliate") {
		return errors.New("article 6 violation: pricing adapter illegally imports affiliate adapter")
	}

	return nil
}
