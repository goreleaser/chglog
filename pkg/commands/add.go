package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/goreleaser/chglog"
)

// nolint: gocognit, funlen, gocritic
func AddCmd(config *viper.Viper) (cmd *cobra.Command) {
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

	cmd.RunE = func(_ *cobra.Command, args []string) (err error) {
		var (
			repoPath string
			gitRepo  *git.Repository
			entries  chglog.ChangeLogEntries
			notes    *chglog.ChangeLogNotes
			data     []byte
		)

		if repoPath, err = os.Getwd(); err != nil {
			return fmt.Errorf("error adding entry: %w", err)
		}

		if len(args) == 1 {
			if repoPath, err = filepath.Abs(args[0]); err != nil {
				return fmt.Errorf("error adding entry: %w", err)
			}
		}

		if entries, err = chglog.Parse(input); err != nil {
			return fmt.Errorf("error adding entry: %w", err)
		}

		if gitRepo, err = chglog.GitRepo(repoPath, true); err != nil {
			return fmt.Errorf("error adding entry: %w", err)
		}

		if headerFile != "" {
			// nolint: gosec, gocritic
			if data, err = os.ReadFile(headerFile); err != nil {
				return fmt.Errorf("error adding entry: %w", err)
			}
			header = string(data)
		}

		if footerFile != "" {
			// nolint: gosec, gocritic
			if data, err = os.ReadFile(footerFile); err != nil {
				return fmt.Errorf("error adding entry: %w", err)
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
			return fmt.Errorf("error adding entry: %w", err)
		}

		if entries, err = chglog.AddEntry(gitRepo, semv, config.GetString("owner"), notes, getDeb(config), entries, config.GetBool("conventional-commits"), config.GetBool("exclude-merge-commits")); err != nil {
			return fmt.Errorf("error adding entry: %w", err)
		}

		if len(entries) == 0 {
			return fmt.Errorf("%w: %s", ErrNoTags, repoPath)
		}

		return entries.Save(output)
	}

	return cmd
}
