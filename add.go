package chglog

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func AddEntry(gitRepo *git.Repository, version *semver.Version, owner string, current ChangeLogEntries, useConventionalCommits bool) (cle ChangeLogEntries, err error) {
	var (
		ref      *plumbing.Reference
		from, to plumbing.Hash
		commits  []*object.Commit
	)

	sort.Sort(current)
	ref, err = gitRepo.Head()
	from = ref.Hash()

	to = plumbing.ZeroHash
	if current.Len() > 0 {
		to, err = GitHashFotTag(gitRepo, current[current.Len()-1].Semver)
	}

	cle = append(cle, current...)
	if commits, err = CommitsBetween(gitRepo, to, from); err != nil {
		return nil, err
	}

	cle = append(cle, CreateEntry(time.Now(), version, owner, commits, useConventionalCommits))
	sort.Sort(sort.Reverse(cle))

	return
}

func processMsg(msg string) string {
	msg = strings.ReplaceAll(strings.ReplaceAll(msg, "\r\n\r\n", "\n\n"), "\r", "")
	msg = regexp.MustCompile(`(?m)(?:^.*Signed-off-by:.*>$)`).ReplaceAllString(msg, "")
	msg = strings.ReplaceAll(strings.Trim(msg, "\n"), "\n\n\n", "\n")
	return msg
}

func CreateEntry(date time.Time, version *semver.Version, owner string, commits []*object.Commit, useConventionalCommits bool) (changelog *ChangeLog) {
	var cc *ConventionalCommit
	changelog = &ChangeLog{
		Semver:   version.String(),
		Date:     date,
		Packager: owner,
	}
	if commits == nil || len(commits) == 0 {
		return
	}
	changelog.Changes = make(ChangeLogChanges, len(commits))

	for idx, c := range commits {
		msg := processMsg(c.Message)
		if useConventionalCommits {
			cc = ParseConventionalCommit(msg)
		}
		changelog.Changes[idx] = &ChangeLogChange{
			Commit:             c.Hash.String(),
			Note:               msg,
			ConventionalCommit: cc,
		}
	}
	return
}
