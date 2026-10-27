package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// Module0052 implements WO-CS-0052: Go monorepo workspace structure and build integrity verifier.
type Module0052 struct {
	baseDir string
}

// NewModule0052 constructs a workspace structure verifier.
func NewModule0052(baseDir string) *Module0052 {
	if baseDir == "" {
		baseDir = "."
	}
	return &Module0052{baseDir: baseDir}
}

// RequiredDirectories specifies the required monorepo architecture folders.
var RequiredDirectories = []string{
	"cmd/api",
	"cmd/worker",
	"internal/platform",
	"internal/platform/core",
	"internal/adapters",
	"internal/domain",
}

// VerifyWorkspace checks that all required Go architectural directories exist.
func (m *Module0052) VerifyWorkspace(ctx context.Context) error {
	for _, dir := range RequiredDirectories {
		fullPath := filepath.Join(m.baseDir, filepath.FromSlash(dir))
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("missing required workspace directory: %s", dir)
			}
			return fmt.Errorf("error inspecting directory %s: %w", dir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path is not a directory: %s", dir)
		}
	}
	return nil
}

