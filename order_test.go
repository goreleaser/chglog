package chglog

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/smartystreets/goconvey/convey"
)

func defSignature() *object.Signature {
	tm, err := time.Parse(time.RFC3339, "2000-01-01T12:00:00+07:00")
	if err != nil {
		tm = time.Now()
	}

	return &object.Signature{
		Name:  "John Doe",
		Email: "John.Doe@example.com",
		When:  tm,
	}
}

func defCommitOptions() *git.CommitOptions {
	return &git.CommitOptions{
		Author:    defSignature(),
		Committer: defSignature(),
	}
}

func newTestRepo() (*git.Repository, error) {
	fs := memfs.New()

	return git.Init(memory.NewStorage(), fs)
}

func TestOrderChangelog(t *testing.T) {
	var (
		gitRepo *git.Repository
		gitTree *git.Worktree
		file    billy.File
		testCLE ChangeLogEntries
	)

	goldcle, err := Parse("./testdata/gold-order-changelog.yml")
	if err != nil {
		t.Fatal(err)
	}

	if gitRepo, err = newTestRepo(); err != nil {
		t.Fatal(err)
	}

	if gitTree, err = gitRepo.Worktree(); err != nil {
		t.Fatal(err)
	}

	if file, err = gitTree.Filesystem.Create("file"); err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	for i := 0; i <= 10; i++ {
		var hash plumbing.Hash

		if _, err = file.Write([]byte(fmt.Sprintf("commit %d\n", i))); err != nil {
			t.Fatal(err)
		}

		if _, err = gitTree.Add("file"); err != nil {
			t.Fatal()
		}

		if hash, err = gitTree.Commit(fmt.Sprintf("commit %d", i), defCommitOptions()); err != nil {
			t.Fatal()
		}

		if _, err = gitRepo.CreateTag(fmt.Sprintf("v0.%d.0", i), hash, nil); err != nil {
			t.Fatal()
		}
	}

	if testCLE, err = InitChangelog(gitRepo, "", nil, nil, false); err != nil {
		t.Fatal(err)
	}

	convey.Convey("Generated entry should be the same as the golden entry", t, func() {
		convey.So(testCLE, convey.ShouldResemble, goldcle)
	})

}
