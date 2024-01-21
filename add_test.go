package chglog

import (
	"log"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/smartystreets/goconvey/convey"
)

func TestAddEntry(t *testing.T) {
	var (
		err     error
		gitRepo *git.Repository
		testCLE ChangeLogEntries
	)
	goldcle, err := Parse("./testdata/gold-add-changelog.yml")
	if err != nil {
		t.Fatal(err)
	}

	if gitRepo, err = GitRepo("./testdata/add-repo", false); err != nil {
		log.Fatal(err)
	}

	testCLE, err = InitChangelog(gitRepo, "", nil, nil, true, false)
	if err != nil {
		t.Fatal(err)
	}
	version, _ := semver.NewVersion("1.0.0-b1+git.123")
	header := `
This is a test
======

header entry
`
	testCLE, err = AddEntry(gitRepo, version, "Dj Gilcrease", createNote(header, ""), nil, testCLE, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(goldcle) != len(testCLE) {
		t.Fatal("differing results")
	}

	// Fix the date since AddEntry uses time.Now
	for i, e := range goldcle {
		testCLE[i].Date = e.Date
	}
	convey.Convey("Generated entry should be the same as the golden entry", t, func() {
		convey.So(testCLE, convey.ShouldResemble, goldcle)
	})
}
