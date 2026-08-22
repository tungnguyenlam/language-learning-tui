package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"deutsch-tui/internal/app"
	"deutsch-tui/internal/content"
	"deutsch-tui/internal/core"
	"deutsch-tui/internal/srs"
	"deutsch-tui/internal/storage/sqlite"
	"deutsch-tui/internal/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	dataDir := flag.String("data-dir", "", "directory for local deutsch-tui data")
	smoke := flag.Bool("smoke", false, "initialize app data and exit")
	testMode := flag.Bool("test-mode", false, "disable dynamic UI elements for testing")
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
	leveledLogger.Info("starting data_dir=%s ai_provider=%s log_level=%s", dir, cfg.AIProvider, cfg.LogLevel)

	store, err := sqlite.Open(filepath.Join(dir, "learning.db"))
	if err != nil {
		leveledLogger.Error("open store: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()

	decks, err := store.Decks(context.Background())
	if err == nil && len(decks) == 0 {
		starter := content.StarterDeck()
		if err := store.UpsertDeck(context.Background(), starter); err != nil {
			leveledLogger.Error("upsert starter deck: %v", err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Also seed dictionary with starter content
		var dictEntries []core.DictionaryEntry
		for _, note := range starter.Notes {
			dictEntries = append(dictEntries, core.DictionaryEntry{
				ID:          note.ID,
				Word:        note.Front,
				Translation: note.Back,
				Tags:        note.Tags,
			})
		}
		if err := store.ImportEntries(context.Background(), dictEntries); err != nil {
			leveledLogger.Error("seed starter dictionary: %v", err)
		}
	}

	importLocalDictCCIfAvailable(store, leveledLogger)

	if *smoke {
		leveledLogger.Info("smoke check complete")
		fmt.Println("deutsch-tui smoke ok")
		return
	}

	scheduler := srs.NewScheduler(leveledLogger)
	program := tea.NewProgram(tui.NewModelWithOptions(store, scheduler, tui.ModelOptions{
		AIProvider:          nil,
		AIProviderName:      cfg.AIProvider,
		DictionaryProvider:  cfg.DictionaryProvider,
		AITemplates:         cfg.AITemplates,
		AISecrets:           secrets,
		TTSProvider:         cfg.TTSProvider,
		TTSVoice:            cfg.TTSVoice,
		TTSCacheDir:         filepath.Join(dir, "audio-cache", "edge-tts"),
		AutoPlayAudio:       cfg.AutoPlayAudio,
		RevealSpeed:         cfg.RevealSpeed,
		StrictNormalization: cfg.StrictNormalization,
		TestMode:            *testMode,
		DataDir:             dir,
		ImportPath:          filepath.Join(dir, "import.tsv"),
		ExportPath:          filepath.Join(dir, "export.tsv"),
		Logger:              leveledLogger,
		OnConfigChange: func(aiProvider string, dictProvider string, tmpls map[string]map[string]string, autoPlayAudio bool, strictNormalization bool, revealSpeed int) {
			pCfg.AIProvider = aiProvider
			pCfg.DictionaryProvider = dictProvider
			pCfg.AITemplates = tmpls
			pCfg.AutoPlayAudio = autoPlayAudio
			pCfg.StrictNormalization = strictNormalization
			pCfg.RevealSpeed = revealSpeed
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

func importLocalDictCCIfAvailable(store *sqlite.Store, leveledLogger *app.LeveledLogger) {
	// Skip local dictionary import in test/E2E environment to avoid performance overhead
	if os.Getenv("PYTEST_CURRENT_TEST") != "" || os.Getenv("DEUTSCH_TUI_BIN") != "" || strings.Contains(filepath.Base(os.Args[0]), "test") {
		return
	}

	ctx := context.Background()
	count, err := store.DictionaryCount(ctx)
	if err != nil {
		leveledLogger.Error("check dictionary count: %v", err)
		return
	}

	// If the database has less than 10,000 dictionary entries, we check if dict.cc files exist
	if count < 10000 {
		files, err := os.ReadDir("local_dict_files")
		if err != nil {
			if !os.IsNotExist(err) {
				leveledLogger.Error("read local_dict_files dir: %v", err)
			}
			return
		}

		var targetZip string
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".zip") {
				targetZip = filepath.Join("local_dict_files", f.Name())
				// Prefer the EN-DE ZIP file (which usually starts with "cngcknmmfb")
				if strings.Contains(f.Name(), "cngcknmmfb") {
					break
				}
			}
		}

		if targetZip != "" {
			fmt.Printf("Found local dictionary zip: %s\n", targetZip)
			fmt.Println("Importing dict.cc offline dictionary entries (this may take a few seconds)...")

			entries, err := content.ParseDictCCZip(targetZip)
			if err != nil {
				leveledLogger.Error("parse dict.cc zip: %v", err)
				fmt.Printf("Error parsing zip: %v\n", err)
				return
			}

			fmt.Printf("Parsed %d entries. Importing into local database...\n", len(entries))

			err = store.ImportEntries(ctx, entries)
			if err != nil {
				leveledLogger.Error("import dict.cc entries: %v", err)
				fmt.Printf("Error importing entries: %v\n", err)
				return
			}

			newCount, _ := store.DictionaryCount(ctx)
			fmt.Printf("Successfully imported %d entries into the local dictionary (Total entries: %d).\n", len(entries), newCount)
		}
	}
}
