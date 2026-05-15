package mcp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeRecordPath_RejectsRelativePath(t *testing.T) {
	_, err := safeRecordPath("relative/output.gif")
	if err == nil || !strings.Contains(err.Error(), "absolute path") {
		t.Fatalf("expected absolute path error, got %v", err)
	}
}

func TestSafeRecordPath_RejectsBadExtension(t *testing.T) {
	_, err := safeRecordPath("/tmp/output.txt")
	if err == nil || !strings.Contains(err.Error(), "unsupported extension") {
		t.Fatalf("expected extension error, got %v", err)
	}
}

func TestSafeRecordPath_RejectsExistingFile(t *testing.T) {
	f := filepath.Join(t.TempDir(), "existing.gif")
	if err := os.WriteFile(f, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := safeRecordPath(f)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected exists error, got %v", err)
	}
}

func TestSafeRecordPath_RejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "real.gif")
	link := filepath.Join(dir, "link.gif")
	if err := os.WriteFile(target, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	_, err := safeRecordPath(link)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink error, got %v", err)
	}
}

func TestSafeRecordPath_AcceptsValidPath(t *testing.T) {
	dir := t.TempDir()
	for _, ext := range []string{".gif", ".webm", ".mp4"} {
		path := filepath.Join(dir, "out"+ext)
		got, err := safeRecordPath(path)
		if err != nil {
			t.Fatalf("ext %s: unexpected error: %v", ext, err)
		}
		if got != path {
			t.Fatalf("ext %s: got %q, want %q", ext, got, path)
		}
	}
}

func TestStreamToFile_WritesAndCaps(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.gif")

	data := strings.NewReader("hello recording data")
	n, err := streamToFile(path, data)
	if err != nil {
		t.Fatalf("streamToFile() error = %v", err)
	}
	if n != 20 {
		t.Fatalf("wrote %d bytes, want 20", n)
	}

	got, _ := os.ReadFile(path)
	if string(got) != "hello recording data" {
		t.Fatalf("file content = %q", got)
	}
}

func TestStreamToFile_RefusesOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.gif")
	if err := os.WriteFile(path, []byte("existing"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := streamToFile(path, strings.NewReader("new data"))
	if err == nil {
		t.Fatal("expected error for existing file")
	}
}
