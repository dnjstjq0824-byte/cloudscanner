package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Module0053 implements WO-CS-0053: Nuxt 3 frontend workspace and SSR hybrid routing verifier.
type Module0053 struct {
	baseDir string
}

// NewModule0053 constructs a frontend workspace verifier.
func NewModule0053(baseDir string) *Module0053 {
	if baseDir == "" {
		baseDir = "."
	}
	return &Module0053{baseDir: baseDir}
}

// RequiredWebFiles specifies the required frontend structure.
var RequiredWebFiles = []string{
	"web/package.json",
	"web/nuxt.config.ts",
	"web/tailwind.config.ts",
	"web/app.vue",
	"web/layouts/public.vue",
	"web/layouts/dashboard.vue",
	"web/layouts/admin.vue",
	"web/pages/index.vue",
}

// VerifyFrontend checks the existence and integrity of the Nuxt 3 frontend workspace.
func (m *Module0053) VerifyFrontend(ctx context.Context) error {
	for _, relFile := range RequiredWebFiles {
		fullPath := filepath.Join(m.baseDir, filepath.FromSlash(relFile))
		if _, err := os.Stat(fullPath); err != nil {
			return fmt.Errorf("missing required frontend file %s: %w", relFile, err)
		}
	}

	// Verify nuxt.config.ts contains noindex for dashboard and admin routes
	nuxtCfgPath := filepath.Join(m.baseDir, "web", "nuxt.config.ts")
	contentBytes, err := os.ReadFile(nuxtCfgPath)
	if err != nil {
		return fmt.Errorf("failed to read nuxt.config.ts: %w", err)
	}
	content := string(contentBytes)

	if !strings.Contains(content, "noindex, nofollow") {
		return fmt.Errorf("nuxt.config.ts must enforce 'noindex, nofollow' for private routes")
	}

	// Verify zero TreeScan references in web package.json
	pkgJsonPath := filepath.Join(m.baseDir, "web", "package.json")
	pkgBytes, err := os.ReadFile(pkgJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read web/package.json: %w", err)
	}
	if strings.Contains(strings.ToLower(string(pkgBytes)), "treescan") {
		return fmt.Errorf("article 1 violation: web/package.json contains treescan coupling")
	}

	return nil
}
