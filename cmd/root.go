package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// Name of this app, injected at build time
var App string

var RootCmd = &cobra.Command{}

func init() {
   viper.SetEnvPrefix(App)
   viper.AutomaticEnv()
}
