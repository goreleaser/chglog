package main

import (
	"strings"

	"github.com/spf13/viper"
)

func setupConfig(cfgFile string) *viper.Viper {
	config := viper.New()
	config.SetEnvPrefix(pkgName)
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
	config.SetDefault("debug", false)
	config.SetConfigFile(cfgFile)
	_ = config.ReadInConfig()
	config.Set("app.name", pkgName)
	config.Set("app.version", version)
	config.Set("app.commit", commit)

	return config
}
