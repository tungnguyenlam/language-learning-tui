package tui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"deutsch-tui/internal/core"

	tea "charm.land/bubbletea/v2"
)

type mockBackupRepo struct {
	*mockRepo
	backupDest string
	restoreSrc string
	backupErr  error
	restoreErr error
}

func (m *mockBackupRepo) Backup(ctx context.Context, destPath string) (core.BackupInfo, error) {
	m.backupDest = destPath
	if m.backupErr != nil {
		return core.BackupInfo{}, m.backupErr
	}
	return core.BackupInfo{Path: destPath, TotalRows: 4, Rows: map[string]int{"cards": 4}}, nil
}

func (m *mockBackupRepo) Restore(ctx context.Context, srcPath string) (core.BackupInfo, error) {
	m.restoreSrc = srcPath
	if m.restoreErr != nil {
		return core.BackupInfo{}, m.restoreErr
	}
	return core.BackupInfo{Path: srcPath, TotalRows: 4, Rows: map[string]int{"cards": 4}}, nil
}

func (m *mockBackupRepo) InspectBackup(ctx context.Context, srcPath string) (core.BackupInfo, error) {
	return core.BackupInfo{Path: srcPath, TotalRows: 4}, nil
}

func TestLatestBackupFilePicksNewest(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{
		"backup-20260101-000000.db",
		"backup-20260815-120000.db",
		"backup-20260304-093000.db",
		"notes.txt",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := latestBackupFile(dir)
	if err != nil {
		t.Fatalf("latestBackupFile: %v", err)
	}
	if filepath.Base(got) != "backup-20260815-120000.db" {
		t.Fatalf("latest = %s, want backup-20260815-120000.db", filepath.Base(got))
	}
}

func TestBackupProgressWritesTimestampedFileViaRepository(t *testing.T) {
	repo := &mockBackupRepo{mockRepo: &mockRepo{}}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		DataDir: t.TempDir(),
	})
	model.activeView = ViewImport

	cmd, handled := model.importScreen.HandleKey(model, tea.KeyPressMsg{Text: "B"})
	if !handled {
		t.Fatal("B was not handled on Import")
	}
	if cmd == nil {
		t.Fatal("backup returned no command")
	}
	msgs := executeCmd(cmd)
	if len(msgs) != 1 {
		t.Fatalf("backup cmd produced %d msgs, want 1: %v", len(msgs), msgs)
	}
	done, ok := msgs[0].(backupDoneMsg)
	if !ok {
		t.Fatalf("backup cmd produced %T, want backupDoneMsg", msgs[0])
	}
	if done.restore {
		t.Fatal("backup reported restore=true")
	}
	if done.info.TotalRows != 4 {
		t.Fatalf("backup rows = %d, want 4", done.info.TotalRows)
	}
	if !strings.HasPrefix(filepath.Base(repo.backupDest), "backup-") || !strings.HasSuffix(repo.backupDest, ".db") {
		t.Fatalf("backup dest = %q, want timestamped backup-*.db", repo.backupDest)
	}
	if filepath.Dir(repo.backupDest) != filepath.Join(model.dataDir, backupsSubdir) {
		t.Fatalf("backup dir = %s, want %s", filepath.Dir(repo.backupDest), filepath.Join(model.dataDir, backupsSubdir))
	}

	updated, follow := model.Update(done)
	model = updated.(*Model)
	if model.lastBackupPath != done.info.Path {
		t.Fatalf("lastBackupPath = %q, want %q", model.lastBackupPath, done.info.Path)
	}
	if !strings.Contains(model.status, "Backed up 4 rows") {
		t.Fatalf("status = %q, want backup confirmation", model.status)
	}
	if follow == nil {
		t.Fatal("backup done returned no status command")
	}
}

func TestBackupProgressRequiresCapableStore(t *testing.T) {
	model := NewModelWithOptions(&mockRepo{}, &mockScheduler{}, ModelOptions{
		DataDir: t.TempDir(),
	})
	model.activeView = ViewImport

	cmd := model.handleBackupProgress()
	if cmd != nil {
		t.Fatal("backup without BackupRepository returned a command")
	}
	if !strings.Contains(model.status, "not available") {
		t.Fatalf("status = %q, want unavailable message", model.status)
	}
}

func TestRestoreProgressConfirmsThenReloads(t *testing.T) {
	repo := &mockBackupRepo{
		mockRepo: &mockRepo{
			decks:    []core.Deck{{ID: "deck-1", Name: "Restored Deck"}},
			dueCards: []core.Card{{ID: "c1", DeckID: "deck-1", Prompt: "Hallo", Answer: "hello"}},
		},
	}
	model := NewModelWithOptions(repo, &mockScheduler{}, ModelOptions{
		DataDir: t.TempDir(),
	})
	model.activeView = ViewImport
	backupPath := filepath.Join(model.dataDir, backupsSubdir, core.BackupFileName(time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)))
	model.lastBackupPath = backupPath

	cmd, handled := model.importScreen.HandleKey(model, tea.KeyPressMsg{Text: "U"})
	if !handled {
		t.Fatal("U was not handled on Import")
	}
	if cmd != nil {
		t.Fatal("restore should wait for confirmation before running")
	}
	if !model.confirmingDelete {
		t.Fatal("restore did not enter confirmation")
	}
	if !strings.Contains(model.status, "RESTORE from") {
		t.Fatalf("status = %q, want restore confirmation", model.status)
	}

	updated, confirmCmd := model.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	model = updated.(*Model)
	if confirmCmd == nil {
		t.Fatal("confirming restore returned no command")
	}
	msgs := executeCmd(confirmCmd)
	if len(msgs) != 1 {
		t.Fatalf("restore cmd produced %d msgs, want 1: %v", len(msgs), msgs)
	}
	done, ok := msgs[0].(backupDoneMsg)
	if !ok || !done.restore {
		t.Fatalf("restore cmd produced %#v, want restore backupDoneMsg", msgs[0])
	}
	if repo.restoreSrc != backupPath {
		t.Fatalf("restore src = %q, want %q", repo.restoreSrc, backupPath)
	}

	updated, _ = model.Update(done)
	model = updated.(*Model)
	if len(model.decks) != 1 || model.decks[0].ID != "deck-1" {
		t.Fatalf("decks after restore = %+v", model.decks)
	}
	if !strings.Contains(model.status, "Restored 4 rows") {
		t.Fatalf("status = %q, want restore confirmation", model.status)
	}
}

func TestImportViewRendersBackupActions(t *testing.T) {
	model := NewModel(&mockRepo{}, &mockScheduler{})
	model.activeView = ViewImport
	model.width, model.height = 120, 40
	view := stripANSI(model.importScreen.Render(model, viewportLayout{}))
	for _, want := range []string{"[B] Backup", "[U] Restore", "Progress backup: none yet"} {
		if !strings.Contains(view, want) {
			t.Fatalf("import view missing %q:\n%s", want, view)
		}
	}
}
