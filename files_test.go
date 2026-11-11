package chglog

import (
	"path/filepath"
	"testing"
)

func TestSaveRoundTripsIndentedNotes(t *testing.T) {
	header := "\n  indented header\nsecond line"
	notes := []string{
		"    Fail closed with malformed allowfrom data (#148)\n\n  * Prepare readme for release\n\n  * Fail closed",
		" leading space\nnext line",
		"\nleading newline\n  indented",
		"plain\n  multi-line",
	}

	changes := make(ChangeLogChanges, 0, len(notes))
	for _, n := range notes {
		changes = append(changes, &ChangeLogChange{Commit: "abc", Note: n})
	}
	entries := ChangeLogEntries{{
		Semver:  "1.0.0",
		Notes:   &ChangeLogNotes{Header: &header},
		Changes: changes,
	}}

	file := filepath.Join(t.TempDir(), "changelog.yml")
	if err := entries.Save(file); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Parse(file)
	if err != nil {
		t.Fatalf("Parse of saved file: %v", err)
	}
	if len(got) != 1 || len(got[0].Changes) != len(notes) {
		t.Fatalf("unexpected entries: %+v", got)
	}
	if *got[0].Notes.Header != header {
		t.Errorf("header = %q, want %q", *got[0].Notes.Header, header)
	}
	for i, want := range notes {
		if got[0].Changes[i].Note != want {
			t.Errorf("note %d = %q, want %q", i, got[0].Changes[i].Note, want)
		}
	}
}
