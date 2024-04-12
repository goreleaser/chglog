package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VersionCmd(config *viper.Viper) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "version",
		Short: "display version info",
	}
	cmd.Run = func(*cobra.Command, []string) {
		version := fmt.Sprintf("%s %s", config.GetString("app.name"), config.GetString("app.version"))
		if config.GetBool("debug") {
			version = fmt.Sprintf("%s+%s", version, config.GetString("app.commit"))
		}
		fmt.Println(version)
	}

	return cmd
}
