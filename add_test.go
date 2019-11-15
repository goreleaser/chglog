package chglog

import (
	"log"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/src-d/go-git.v4"
)

func TestAddEntry(t *testing.T) {
	var (
		err     error
		gitRepo *git.Repository
		testCLE ChangeLogEntries
	)
	if gitRepo, err = GitRepo("./testdata/add-repo", false); err != nil {
		log.Fatal(err)
	}

	testCLE, err = InitChangelog(gitRepo, "", nil, nil, true)
	if err != nil {
		t.Error(err)
		return
	}
	version, _ := semver.NewVersion("1.0.0-b1+git.123")
	header := `
This is a test
======

header entry
`
	testCLE, err = AddEntry(gitRepo, version, "Dj Gilcrease", createNote(header, ""), nil, testCLE, true)
	if err != nil {
		t.Error(err)
		return
	}
	// Fix the date since AddEntry uses time.Now
	testCLE[0].Date, _ = time.Parse(time.RFC3339Nano, "2019-10-18T18:17:57.934767812-07:00")
	goldcle, err := Parse("./testdata/gold-add-changelog.yml")
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(goldcle, testCLE); diff != "" {
		t.Errorf("ChangeLogEntries mismatch (+got -want):\n%s", diff)
	}
}
