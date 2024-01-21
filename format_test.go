package chglog

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-cmp/cmp"
)

// nolint: gochecknoglobals,gocritic
var formats = map[string]string{"rpm": rpmTpl, "deb": debTpl, "release": releaseTpl, "repo": repoTpl}

func TestFormatChangelog(t *testing.T) {
	var (
		err     error
		gitRepo *git.Repository
		testCLE ChangeLogEntries
	)
	t.Parallel()
	if gitRepo, err = GitRepo("./testdata/init-repo", false); err != nil {
		log.Fatal(err)
	}

	if testCLE, err = InitChangelog(gitRepo, "", nil, nil, true, false); err != nil {
		t.Error(err)

		return
	}

	for tmplType, tmplData := range formats {
		tmplData := tmplData
		pkg := PackageChangeLog{fmt.Sprintf("TestFormatChangelog-%s", tmplType), testCLE}
		t.Run(tmplType, func(t *testing.T) {
			t.Parallel()
			accept(t, tmplData, pkg)
		})
	}
}

func accept(t *testing.T, tmplData string, pkg PackageChangeLog) {
	t.Helper()

	tpl, err := LoadTemplateData(tmplData)
	if err != nil {
		t.Error(err)

		return
	}

	testdata, err := FormatChangelog(&pkg, tpl)
	if err != nil {
		t.Error(err)

		return
	}

	golddata, _ := os.ReadFile(fmt.Sprintf("./testdata/%s", pkg.Name))
	if diff := cmp.Diff(string(golddata), testdata); diff != "" {
		t.Errorf("FormatChangelog mismatch (+got -want):\n%s", diff)
	}
}
