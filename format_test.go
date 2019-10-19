package chglog

import (
	"fmt"
	"log"
	"testing"

	"gopkg.in/src-d/go-git.v4"
)

// nolint: gochecknoglobals
var formats = map[string]string{"rpm": rpmTpl, "deb": debTpl, "repo": repoTpl}

func TestFormatChangelog(t *testing.T) {
	var (
		err     error
		gitRepo *git.Repository
		testCLE ChangeLogEntries
	)
	if gitRepo, err = GitRepo("./testdata/init-repo"); err != nil {
		log.Fatal(err)
	}

	if testCLE, err = InitChangelog(gitRepo, "", nil, nil, true); err != nil {
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
	if tpl, err := LoadTemplateData(tmplData); err != nil {
		t.Error(err)
		return
	} else if testdata, err := FormatChangelog(&pkg, tpl); err != nil {
		t.Error(err)
		return
	} else {
		fmt.Printf("\n\n===============\n%s\n===============\n", testdata)
	}
}
