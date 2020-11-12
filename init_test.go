package chglog

import (
	"log"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/smartystreets/goconvey/convey"
)

func TestInitChangelog(t *testing.T) {
	var (
		err     error
		gitRepo *git.Repository
		testCLE ChangeLogEntries
	)

	goldcle, err := Parse("./testdata/gold-init-changelog.yml")
	if err != nil {
		t.Fatal(err)
	}
	if gitRepo, err = GitRepo("./testdata/init-repo", false); err != nil {
		log.Fatal(err)
	}

	testCLE, err = InitChangelog(gitRepo, "", nil, nil, true)
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
