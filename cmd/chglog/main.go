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
	debug := false
	config := setupConfig(cfgFile)

	cmdRoot.PersistentFlags().BoolVarP(
		&debug,
		"debug",
		"",
		debug,
		``)
	config.BindPFlag("debug", cmdRoot.PersistentFlags().Lookup("debug"))

	cmdRoot.PersistentFlags().StringVarP(
		&cfgFile,
		"config-file",
		"c",
		cfgFile,
		``)
	config.BindPFlag("config-file", cmdRoot.PersistentFlags().Lookup("config-file"))

	cmdRoot.PersistentPreRun = func(c *cobra.Command, args []string) {
		if cfgFile == config.GetString("config-file") {
			return
		}

		config.SetConfigFile(cfgFile)
		config.Set("config-file", cfgFile)
		config.ReadInConfig()
	}

	cmds := commands.AllCommands(config)
	cmdRoot.AddCommand(cmds...)

	if err := cmdRoot.Execute(); err != nil {
		// nolint: gomnd, gocritic
		os.Exit(127)
	}
}
