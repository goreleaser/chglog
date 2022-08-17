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
	config, err := setupConfig(cfgFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(127)
	}

	cmdRoot.PersistentFlags().StringVarP(
		&cfgFile,
		"config-file",
		"c",
		cfgFile,
		``)
	_ = config.BindPFlag("config-file", cmdRoot.PersistentFlags().Lookup("config-file"))

	cmdRoot.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		if cfgFile == config.GetString("config-file") {
			return nil
		}

		config.SetConfigFile(cfgFile)
		config.Set("config-file", cfgFile)
		return config.ReadInConfig()
	}

	cmds := commands.AllCommands(config)
	cmdRoot.AddCommand(cmds...)

	if err := cmdRoot.Execute(); err != nil {
		// nolint: gomnd, gocritic
		os.Exit(127)
	}
}
