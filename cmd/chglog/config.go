package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func setupConfig(cfgFile string) (*viper.Viper, error) {
	config := viper.New()
	config.SetEnvPrefix(pkgName)
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
	config.SetDefault("debug", false)
	config.SetConfigFile(cfgFile)
	if err := config.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}
	config.Set("app.name", pkgName)
	config.Set("app.version", version)
	config.Set("app.commit", commit)

	return config, nil
}
