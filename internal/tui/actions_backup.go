package tui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

const backupsSubdir = "backups"

type backupDoneMsg struct {
	restore bool
	info    core.BackupInfo
	decks   []core.Deck
	cards   []core.Card
}

var errNoBackups = errors.New("no progress backups found")

type latestBackupPathMsg struct {
	id      int
	path    string
	err     error
	restore bool
}

func (m *Model) backupDir() string {
	if strings.TrimSpace(m.dataDir) == "" {
		return ""
	}
	return filepath.Join(m.dataDir, backupsSubdir)
}

func (m *Model) backupRepository() (core.BackupRepository, bool) {
	repo, ok := m.repo.(core.BackupRepository)
	return repo, ok
}

func latestBackupFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "backup-") && strings.HasSuffix(name, ".db") {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return "", fmt.Errorf("%w in %s", errNoBackups, dir)
	}
	sort.Strings(names)
	return filepath.Join(dir, names[len(names)-1]), nil
}

func (m *Model) loadLatestBackupPath(restore bool) tea.Cmd {
	dir := m.backupDir()
	m.backupPathLoadID++
	id := m.backupPathLoadID
	return func() tea.Msg {
		if dir == "" {
			return latestBackupPathMsg{id: id, restore: restore, err: errNoBackups}
		}
		path, err := latestBackupFile(dir)
		return latestBackupPathMsg{id: id, path: path, err: err, restore: restore}
	}
}

func (m *Model) handleBackupProgress() tea.Cmd {
	if _, ok := m.backupRepository(); !ok {
		m.status = "Progress backup is not available for this store."
		return nil
	}
	dir := m.backupDir()
	if dir == "" {
		m.status = "No data directory configured; cannot write a backup."
		return nil
	}
	dest := filepath.Join(dir, core.BackupFileName(time.Now()))
	m.status = "Backing up progress..."
	return m.executeBackup(dest)
}

func (m *Model) executeBackup(dest string) tea.Cmd {
	repo, ok := m.backupRepository()
	if !ok {
		m.status = "Progress backup is not available for this store."
		return nil
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		info, err := repo.Backup(ctx, dest)
		if err != nil {
			return fmt.Errorf("backup failed: %w", err)
		}
		return backupDoneMsg{info: info}
	}
}

func (m *Model) handleRestoreProgress() tea.Cmd {
	if _, ok := m.backupRepository(); !ok {
		m.status = "Progress restore is not available for this store."
		return nil
	}
	path := strings.TrimSpace(m.lastBackupPath)
	if path == "" {
		m.status = "Finding latest progress backup..."
		return m.loadLatestBackupPath(true)
	}
	return m.confirmRestoreProgress(path)
}

func (m *Model) confirmRestoreProgress(path string) tea.Cmd {
	m.confirmingDelete = true
	m.deleteAction = func() tea.Cmd { return m.executeRestore(path) }
	m.status = fmt.Sprintf("RESTORE from %s? Replaces decks and progress. (y/n)", filepath.Base(path))
	return nil
}

func (m *Model) handleLatestBackupPathMsg(msg latestBackupPathMsg) tea.Cmd {
	if msg.id != m.backupPathLoadID || m.activeView != ViewImport {
		return nil
	}
	if msg.err != nil {
		m.lastBackupPath = ""
		if errors.Is(msg.err, errNoBackups) || errors.Is(msg.err, os.ErrNotExist) {
			if msg.restore {
				m.status = "No progress backup found. Press B to create one."
			}
			return nil
		}
		return func() tea.Msg { return fmt.Errorf("discover progress backup: %w", msg.err) }
	}
	m.lastBackupPath = msg.path
	if msg.restore {
		return m.confirmRestoreProgress(msg.path)
	}
	return nil
}

func (m *Model) executeRestore(src string) tea.Cmd {
	repo, ok := m.backupRepository()
	if !ok {
		m.status = "Progress restore is not available for this store."
		return nil
	}
	m.status = "Restoring progress..."
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		info, err := repo.Restore(ctx, src)
		if err != nil {
			return fmt.Errorf("restore failed: %w", err)
		}
		decks, err := m.repo.Decks(ctx)
		if err != nil {
			return err
		}
		cards, err := m.repo.DueCards(ctx, time.Now(), 0)
		if err != nil {
			return err
		}
		return backupDoneMsg{restore: true, info: info, decks: decks, cards: cards}
	}
}
