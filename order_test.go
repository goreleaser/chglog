package chglog

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	. "github.com/smartystreets/goconvey/convey"
)

type testRepo struct {
	Git    *git.Repository
	Source *git.Worktree
	seqno  int
}

func newTestRepo() *testRepo {
	var (
		repo *git.Repository
		tree *git.Worktree
		err  error
	)

	fs := memfs.New()

	if repo, err = git.Init(memory.NewStorage(), fs); err != nil {
		log.Fatal(err)
	}

	if tree, err = repo.Worktree(); err != nil {
		log.Fatal(err)
	}

	return &testRepo{
		Git:    repo,
		Source: tree,
	}
}

// modifyAndCommit creates the file if it does not exist, appends a
// change, commits the file, and returns the hash of the commit.
func (r *testRepo) modifyAndCommit(filename string, opts *git.CommitOptions) plumbing.Hash {
	var (
		hash plumbing.Hash
		err  error
		file billy.File
	)

	if file, err = r.Source.Filesystem.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err = file.Write([]byte(fmt.Sprintf("commit %d\n", r.seqno))); err != nil {
		log.Fatal(err)
	}

	if _, err = r.Source.Add(filename); err != nil {
		log.Fatal(err)
	}

	if hash, err = r.Source.Commit(fmt.Sprintf("commit %d", r.seqno), opts); err != nil {
		log.Fatal(err)
	}

	r.seqno++

	return hash
}

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

func TestOrderChangelog(t *testing.T) {
	goldCLE, err := Parse("./testdata/gold-order-changelog.yml")
	if err != nil {
		t.Fatal(err)
	}

	repo := newTestRepo()

	for i := 0; i <= 10; i++ {
		hash := repo.modifyAndCommit("file", defCommitOptions())

		if _, err = repo.Git.CreateTag(fmt.Sprintf("v0.%d.0", i), hash, nil); err != nil {
			t.Fatal(err)
		}
	}

	testCLE, err := InitChangelog(repo.Git, "", nil, nil, false)
	if err != nil {
		t.Fatal(err)
	}

	Convey("Generated entry should be the same as the golden entry", t, func() {
		So(testCLE, ShouldResemble, goldCLE)
	})

}

func TestSemverTag(t *testing.T) {
	repo := newTestRepo()
	tag := "1.0.0"

	Convey("Semver tags should be parsed", t, func() {
		hash := repo.modifyAndCommit("file", defCommitOptions())

		if _, err := repo.Git.CreateTag(tag, hash, nil); err != nil {
			t.Fatal(err)
		}

		cle, err := InitChangelog(repo.Git, "", nil, nil, false)
		if err != nil {
			t.Fatal(err)
		}

		So(cle, ShouldHaveLength, 1)
		So(cle[0].Semver, ShouldEqual, tag)
	})

	Convey("Not Semver tags should be ignored", t, func() {
		hash := repo.modifyAndCommit("file", defCommitOptions())

		if _, err := repo.Git.CreateTag("text", hash, nil); err != nil {
			t.Fatal(err)
		}

		cle, err := InitChangelog(repo.Git, "", nil, nil, false)
		if err != nil {
			t.Fatal(err)
		}

		So(cle, ShouldHaveLength, 1)
		So(cle[0].Semver, ShouldEqual, tag)
	})
}
