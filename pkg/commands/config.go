package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setupConfigCmd(config *viper.Viper) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "config",
		Short: "save config data",
	}

	cmd, config = commonFlags(cmd, config)
	cmd.PersistentPreRun = func(c *cobra.Command, args []string) {
		cmd.Parent().PersistentPreRun(c, args)
	}

	cmd.RunE = func(c *cobra.Command, args []string) (err error) {
		// Filter some config settings
		cfgMap := config.AllSettings()
		delete(cfgMap, "app")
		delete(cfgMap, "config-file")
		v := viper.New()
		v.MergeConfigMap(cfgMap)

		return v.WriteConfigAs(config.GetString("config-file"))
	}

	return cmd
}
