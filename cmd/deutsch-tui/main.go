package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	pCfg := &cfg
	secrets, err := app.LoadSecrets(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load secrets:", err)
		os.Exit(1)
	}
	pSecrets := &secrets
	logFile, logger, err := app.OpenLog(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Create leveled logger based on config
	level := app.ParseLogLevel(cfg.LogLevel)
	leveledLogger := app.NewLeveledLogger(logger, level)
	leveledLogger.Info("starting data_dir=%s theme=%s keymap=%s ai_provider=%s log_level=%s", dir, cfg.Theme, cfg.Keymap, cfg.AIProvider, cfg.LogLevel)

	store, err := sqlite.Open(filepath.Join(dir, "learning.db"))
	if err != nil {
		leveledLogger.Error("open store: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()

	decks, err := store.Decks(context.Background())
	if err == nil && len(decks) == 0 {
		if err := store.UpsertDeck(context.Background(), content.StarterDeck()); err != nil {
			leveledLogger.Error("upsert starter deck: %v", err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if *smoke {
		leveledLogger.Info("smoke check complete")
		fmt.Println("deutsch-tui smoke ok")
		return
	}

	scheduler := srs.NewScheduler(leveledLogger)
	program := tea.NewProgram(tui.NewModelWithOptions(store, scheduler, tui.ModelOptions{
		Theme:               cfg.Theme,
		AIProvider:          nil,
		AIProviderName:      cfg.AIProvider,
		AITemplates:         cfg.AITemplates,
		AISecrets:           secrets,
		AutoPlayAudio:       cfg.AutoPlayAudio,
		StrictNormalization: cfg.StrictNormalization,
		ImportPath:          filepath.Join(dir, "import.tsv"),
		ExportPath:          filepath.Join(dir, "export.tsv"),
		Logger:              leveledLogger, // Pass the logger
		OnConfigChange: func(theme string, name string, tmpls map[string]map[string]string, autoPlayAudio bool, strictNormalization bool) {
			pCfg.Theme = theme
			pCfg.AIProvider = name
			pCfg.AITemplates = tmpls
			pCfg.AutoPlayAudio = autoPlayAudio
			pCfg.StrictNormalization = strictNormalization
			if err := app.SaveConfig(dir, *pCfg); err != nil {
				leveledLogger.Error("save config: %v", err)
			}
		},
		OnSecretsChange: func(s app.Secrets) {
			*pSecrets = s
			if err := app.SaveSecrets(dir, *pSecrets); err != nil {
				leveledLogger.Error("save secrets: %v", err)
			}
		},
	}))
	if _, err := program.Run(); err != nil {
		leveledLogger.Error("program run: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
