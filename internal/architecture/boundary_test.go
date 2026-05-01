package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageBoundaries(t *testing.T) {
	rules := []struct {
		dir       string
		forbidden []string
	}{
		{dir: "internal/core", forbidden: []string{"deutsch-tui/internal/tui", "deutsch-tui/internal/storage", "deutsch-tui/internal/ai", "database/sql", "charm.land/"}},
		{dir: "internal/storage", forbidden: []string{"deutsch-tui/internal/tui", "deutsch-tui/internal/ai", "charm.land/"}},
		{dir: "internal/content", forbidden: []string{"deutsch-tui/internal/tui", "deutsch-tui/internal/storage", "charm.land/"}},
		{dir: "internal/ai", forbidden: []string{"deutsch-tui/internal/tui", "deutsch-tui/internal/storage", "database/sql", "charm.land/"}},
	}

	root := repoRoot(t)
	for _, rule := range rules {
		rule := rule
		t.Run(rule.dir, func(t *testing.T) {
			checkForbiddenImports(t, filepath.Join(root, rule.dir), rule.forbidden)
		})
	}
}

func checkForbiddenImports(t *testing.T, dir string, forbidden []string) {
	t.Helper()
	fset := token.NewFileSet()
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, imported := range file.Imports {
			value := strings.Trim(imported.Path.Value, `"`)
			for _, prefix := range forbidden {
				if strings.HasPrefix(value, prefix) {
					t.Fatalf("%s imports forbidden package %q", path, value)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		next := filepath.Dir(dir)
		if next == dir {
			t.Fatal("repo root not found")
		}
		dir = next
	}
}
