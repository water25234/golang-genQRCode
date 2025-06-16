package cmd

import "testing"

func TestCheckFlags(t *testing.T) {
	// success case
	readfile = "file.txt"
	folder = "dir"
	fileExt = ".png"
	fs, err := checkFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.readfile != readfile || fs.folder != folder || fs.fileExt != fileExt {
		t.Errorf("returned struct mismatch")
	}

	// missing readfile
	readfile = ""
	_, err = checkFlags()
	if err == nil {
		t.Errorf("expected error when readfile empty")
	}

	// restore for other tests
	readfile = "file.txt"
	folder = "dir"
	fileExt = ".png"
}
