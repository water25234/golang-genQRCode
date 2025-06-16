package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPinCodeInfo(t *testing.T) {
	fs := &flags{}
	tests := []struct {
		name      string
		input     []string
		wantName  string
		wantPin   string
		wantError bool
	}{
		{"single", []string{"123"}, "123", "123", false},
		{"pair", []string{"abc", "789"}, "abc", "789", false},
		{"invalid", []string{"a", "b", "c"}, "", "", true},
	}

	for _, tt := range tests {
		gotName, gotPin, err := fs.pinCodeInfo(tt.input)
		if tt.wantError {
			if err == nil {
				t.Errorf("%s: expected error", tt.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}
		if gotName != tt.wantName || gotPin != tt.wantPin {
			t.Errorf("%s: want %s/%s got %s/%s", tt.name, tt.wantName, tt.wantPin, gotName, gotPin)
		}
	}
}

func TestFileSize(t *testing.T) {
	fs := &flags{}
	dir := t.TempDir()
	p := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(p, []byte("abc"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	size, err := fs.fileSize(p)
	if err != nil {
		t.Fatalf("fileSize error: %v", err)
	}
	if size != 3 {
		t.Errorf("expected size 3 got %d", size)
	}
}
