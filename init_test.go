package chglog

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInitChangelog(t *testing.T) {
	gitRepo, err := GitRepo("./testdata/repo")
	if err != nil {
		t.Error(err)
	}

	testcle, err := InitChangelog(gitRepo, true)
	if err != nil {
		t.Error(err)
	}
	goldcle, err := Parse("./testdata/gold-changelog.yml")

	if diff := cmp.Diff(goldcle, testcle); diff != "" {
		t.Errorf("ChangeLogEntries mismatch (-got +want):\n%s", diff)
	}
}
