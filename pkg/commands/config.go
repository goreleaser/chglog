package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ConfigCmd(config *viper.Viper) (cmd *cobra.Command) {
	var pkg string
	cmd = &cobra.Command{
		Use:   "config",
		Short: "save config data",
	}

	cmd, config = commonFlags(cmd, config)
	cmd.Flags().StringVarP(
		&pkg,
		"package-name",
		"p",
		"",
		"package name to use in formatting")

	cmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		return config.BindPFlag("package-name", cmd.Flags().Lookup("package-name"))
	}

	cmd.RunE = func(*cobra.Command, []string) error {
		// Filter some config settings
		cfgMap := config.AllSettings()
		delete(cfgMap, "app")
		delete(cfgMap, "config-file")
		v := viper.New()
		if err := v.MergeConfigMap(cfgMap); err != nil {
			return err
		}

		return v.WriteConfigAs(config.GetString("config-file"))
	}

	return cmd
}
