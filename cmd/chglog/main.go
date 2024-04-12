// Package main contains the main nfpm cli source code.
package main

import (
	"fmt"
	"os"
	"path"

	"github.com/goreleaser/chglog/pkg/commands"
	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals, gocritic
var (
	pkgName = "chglog"
	version = "v0.0.0"
	commit  = "local"
)

func main() {
	cmdRoot := &cobra.Command{
		Use:          pkgName,
		Short:        "Changelog generator",
		SilenceUsage: true,
	}
	cwd, _ := os.Getwd()
	cfgFile := path.Join(cwd, fmt.Sprintf(".%s.yml", pkgName))
	config := setupConfig(cfgFile)

	cmdRoot.PersistentFlags().StringVarP(
		&cfgFile,
		"config-file",
		"c",
		cfgFile,
		``)
	_ = config.BindPFlag("config-file", cmdRoot.PersistentFlags().Lookup("config-file"))

	cmdRoot.PersistentPreRunE = func(*cobra.Command, []string) error {
		if cfgFile == config.GetString("config-file") {
			return nil
		}

		config.SetConfigFile(cfgFile)
		config.Set("config-file", cfgFile)
		return config.ReadInConfig()
	}

	cmdRoot.AddCommand(
		commands.AddCmd(config),
		commands.InitCmd(config),
		commands.VersionCmd(config),
		commands.ConfigCmd(config),
		commands.FormatCmd(config),
	)

	if err := cmdRoot.Execute(); err != nil {
		// nolint: gomnd, gocritic
		os.Exit(127)
	}
}
