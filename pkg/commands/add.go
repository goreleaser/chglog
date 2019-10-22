package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"

	"github.com/goreleaser/chglog"
)

// nolint: gocognit, funlen
func setupAddCmd(config *viper.Viper) (cmd *cobra.Command) {
	var (
		input,
		output,
		version,
		header,
		footer,
		headerFile,
		footerFile string
		semv *semver.Version
	)
	cmd, config = commonFlags(&cobra.Command{
		Use:   "add [PATH]",
		Short: "add new changelog entry for [PATH]",
		Args:  cobra.MaximumNArgs(1),
	}, config)

	cmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"changelog.yml",
		"file to save the updated changelog to")
	cmd.Flags().StringVarP(
		&input,
		"input",
		"i",
		"changelog.yml",
		"starting changelog.yml file")
	cmd.Flags().StringVarP(
		&version,
		"version",
		"",
		"",
		"Version for this entry")
	cmd.Flags().StringVarP(
		&header,
		"header",
		"",
		"",
		"Header note for this entry")
	cmd.Flags().StringVarP(
		&footer,
		"footer",
		"",
		"",
		"Footer note for this entry")
	cmd.Flags().StringVarP(
		&headerFile,
		"header-file",
		"",
		"",
		"Header note for this entry")
	cmd.Flags().StringVarP(
		&footerFile,
		"footer-file",
		"",
		"",
		"Footer note for this entry")

	cmd.PersistentPreRun = func(c *cobra.Command, args []string) {
		cmd.Parent().PersistentPreRun(c, args)
	}

	cmd.RunE = func(c *cobra.Command, args []string) (err error) {
		var (
			repoPath string
			gitRepo  *git.Repository
			entries  chglog.ChangeLogEntries
			notes    *chglog.ChangeLogNotes
			data     []byte
		)

		if repoPath, err = os.Getwd(); err != nil {
			return err
		}

		if len(args) == 1 {
			if repoPath, err = filepath.Abs(args[0]); err != nil {
				return err
			}
		}

		if entries, err = chglog.Parse(input); err != nil {
			return err
		}

		if gitRepo, err = chglog.GitRepo(repoPath); err != nil {
			return err
		}

		if headerFile != "" {
			if data, err = ioutil.ReadFile(headerFile); err != nil { // nolint: gosec
				return err
			}
			header = string(data)
		}

		if footerFile != "" {
			if data, err = ioutil.ReadFile(footerFile); err != nil { // nolint: gosec
				return err
			}
			footer = string(data)
		}

		if header != "" || footer != "" {
			notes = &chglog.ChangeLogNotes{}
			if header != "" {
				header = strings.ReplaceAll(header, "\\n", "\n")
				notes.Header = &header
			}
			if footer != "" {
				notes.Footer = &footer
			}
		}

		if semv, err = semver.NewVersion(version); err != nil {
			return err
		}

		if entries, err = chglog.AddEntry(gitRepo, semv, config.GetString("owner"), notes, getDeb(config), entries, config.GetBool("conventional-commits")); err != nil {
			return err
		}

		if len(entries) == 0 {
			return fmt.Errorf("%s does not have any versioned releases. `git tag` should return semver formated tags", repoPath)
		}

		return entries.Save(output)
	}

	return cmd
}
