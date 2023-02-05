// Package commands contain the commands for the cli
package commands

import (
	"github.com/goreleaser/chglog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func commonFlags(cmd *cobra.Command, config *viper.Viper) (*cobra.Command, *viper.Viper) {
	var (
		urgency, owner      string
		distribution        []string
		conventionalCommits bool
	)
	cmd.Flags().BoolVarP(
		&conventionalCommits,
		"conventional-commits",
		"",
		conventionalCommits,
		`Use conventional commits parsing`)
	cmd.Flags().StringVarP(
		&owner,
		"owner",
		"",
		owner,
		`set package owner`)
	cmd.Flags().StringVarP(
		&urgency,
		"deb-urgency",
		"",
		urgency,
		`set debian urgency for`)
	cmd.Flags().StringSliceVarP(
		&distribution,
		"deb-distribution",
		"",
		distribution,
		`set debian distributions for`)

	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		if err := config.BindPFlag("conventional-commits", cmd.Flags().Lookup("conventional-commits")); err != nil {
			return err
		}
		if err := config.BindPFlag("owner", cmd.Flags().Lookup("owner")); err != nil {
			return err
		}
		if err := config.BindPFlag("deb.urgency", cmd.Flags().Lookup("deb-urgency")); err != nil {
			return err
		}
		return config.BindPFlag("deb.distribution", cmd.Flags().Lookup("deb-distribution"))
	}

	return cmd, config
}

func getDeb(config *viper.Viper) (deb *chglog.ChangelogDeb) {
	var (
		urgency       string
		distributions []string
	)
	urgency = config.GetString("deb.urgency")
	distributions = config.GetStringSlice("deb.distribution")
	if len(distributions) > 0 && urgency != "" {
		deb = &chglog.ChangelogDeb{
			Urgency:       urgency,
			Distributions: distributions,
		}
	}

	return deb
}
