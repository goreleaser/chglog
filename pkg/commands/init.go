package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"

	"github.com/goreleaser/chglog"
)

func setupInitCmd(config *viper.Viper) (cmd *cobra.Command) {
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

	cmd.PersistentPreRun = func(c *cobra.Command, args []string) {
		cmd.Parent().PersistentPreRun(c, args)
	}

	cmd.RunE = func(c *cobra.Command, args []string) (err error) {
		var (
			repoPath string
			gitRepo  *git.Repository
			entries  chglog.ChangeLogEntries
		)
		if repoPath, err = os.Getwd(); err != nil {
			return err
		}

		if len(args) == 1 {
			if repoPath, err = filepath.Abs(args[0]); err != nil {
				return err
			}
		}

		if gitRepo, err = chglog.GitRepo(repoPath, true); err != nil {
			return err
		}

		if entries, err = chglog.InitChangelog(gitRepo, config.GetString("owner"), nil, getDeb(config), config.GetBool("conventional-commits")); err != nil {
			return err
		}

		if len(entries) == 0 {
			return fmt.Errorf("%s does not have any versioned releases. `git tag` should return semver formated tags", repoPath)
		}

		return entries.Save(output)
	}

	return cmd
}
