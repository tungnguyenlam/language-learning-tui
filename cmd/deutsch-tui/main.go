package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"deutsch-tui/internal/ai"
	"deutsch-tui/internal/app"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/srs"
	"deutsch-tui/internal/storage/sqlite"
	"deutsch-tui/internal/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	dataDir := flag.String("data-dir", "", "directory for local deutsch-tui data")
	smoke := flag.Bool("smoke", false, "initialize app data and exit")
	flag.Parse()

	dir, err := app.ResolveDataDir(*dataDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := app.EnsureDataDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cfg, err := app.LoadOrCreateConfig(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logFile, logger, err := app.OpenLog(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer logFile.Close()
	logger.Printf("starting data_dir=%s theme=%s keymap=%s ai_provider=%s", dir, cfg.Theme, cfg.Keymap, cfg.AIProvider)

	store, err := sqlite.Open(filepath.Join(dir, "learning.db"))
	if err != nil {
		logger.Printf("open store: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()

	if err := store.UpsertDeck(context.Background(), content.StarterDeck()); err != nil {
		logger.Printf("upsert starter deck: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *smoke {
		logger.Print("smoke check complete")
		fmt.Println("deutsch-tui smoke ok")
		return
	}

	scheduler := srs.NewScheduler()
	var provider ai.Provider
	switch cfg.AIProvider {
	case "template":
		provider = ai.TemplateProvider{Templates: cfg.AITemplates}
	case "offline":
		provider = ai.OfflineProvider{}
	default:
		provider = ai.OfflineProvider{}
	}
	program := tea.NewProgram(tui.NewModelWithOptions(store, scheduler, tui.ModelOptions{
		AIProvider:     provider,
		AIProviderName: cfg.AIProvider,
		AITemplates:    cfg.AITemplates,
		AutoPlayAudio:  cfg.AutoPlayAudio,
		ImportPath:     filepath.Join(dir, "import.tsv"),
		ExportPath:     filepath.Join(dir, "export.tsv"),
		OnConfigChange: func(name string, tmpls map[string]string, autoPlayAudio bool) {
			cfg.AIProvider = name
			cfg.AITemplates = tmpls
			cfg.AutoPlayAudio = autoPlayAudio
			if err := app.SaveConfig(dir, cfg); err != nil {
				logger.Printf("save config: %v", err)
			}
		},
	}))
	if _, err := program.Run(); err != nil {
		logger.Printf("program run: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
