package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/goreleaser/chglog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ErrNoTags happens when a repository has no tags.
var ErrNoTags = errors.New("no versioned releases found, check the output of `git tag`")

func InitCmd(config *viper.Viper) (cmd *cobra.Command) {
	var output string
	cmd = &cobra.Command{
		Use:   "init [PATH]",
		Short: "create a new changelog file for [PATH]",
		Args:  cobra.MaximumNArgs(1),
	}

	cmd, config = commonFlags(cmd, config)
	cmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"changelog.yml",
		"file to save the new changelog to")

	cmd.RunE = func(_ *cobra.Command, args []string) (err error) {
		var (
			repoPath string
			gitRepo  *git.Repository
			entries  chglog.ChangeLogEntries
		)
		if repoPath, err = os.Getwd(); err != nil {
			return fmt.Errorf("error initialzing change log: %w", err)
		}

		if len(args) == 1 {
			if repoPath, err = filepath.Abs(args[0]); err != nil {
				return fmt.Errorf("error initialzing change log: %w", err)
			}
		}

		if gitRepo, err = chglog.GitRepo(repoPath, true); err != nil {
			return fmt.Errorf("error initialzing change log: %w", err)
		}

		if entries, err = chglog.InitChangelog(gitRepo, config.GetString("owner"), nil, getDeb(config), config.GetBool("conventional-commits"), config.GetBool("exclude-merge-commits")); err != nil {
			return fmt.Errorf("error initialzing change log: %w", err)
		}

		if len(entries) == 0 {
			return fmt.Errorf("%w: %s", ErrNoTags, repoPath)
		}

		if err := entries.Save(output); err != nil {
			return err
		}

		fmt.Println("created:", output)
		return nil
	}

	return cmd
}
