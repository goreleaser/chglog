package chglog

import (
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/src-d/go-git.v4"
)

func TestInitChangelog(t *testing.T) {
	var (
		err     error
		gitRepo *git.Repository
		testCLE ChangeLogEntries
	)
	if gitRepo, err = GitRepo("./testdata/init-repo"); err != nil {
		log.Fatal(err)
	}

	testCLE, err = InitChangelog(gitRepo, "", nil, nil, true)
	if err != nil {
		t.Error(err)
		return
	}

	goldcle, err := Parse("./testdata/gold-init-changelog.yml")
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(goldcle, testCLE); diff != "" {
		t.Errorf("ChangeLogEntries mismatch (-got +want):\n%s", diff)
	}
}
