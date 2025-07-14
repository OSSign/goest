package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ossign/goest/lib"
	"github.com/ossign/goest/lib/constants"
	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "goest",
	Short: "Go SignTool - A cross-platform tool for code signing",
	Long:  "Go SignTool is a cross-platform tool for code signing based on Relic from sassoftware.",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.PersistentFlags().GetString("config")
		fmt.Println("Config file:", configFile)
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if configFile != "" {
			fmt.Println("Using configuration file:", configFile)

			Viper.SetConfigFile(configFile)
			Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			Viper.SetEnvPrefix("GOEST")
			Viper.AutomaticEnv()

			if err := Viper.ReadInConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Using config file: %s\n", Viper.ConfigFileUsed())

			if err := Viper.Unmarshal(&GlobalConfig); err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshalling config: %v\n", err)
				os.Exit(1)
			}

			lib.RegisterSigners(defaultSigners)
		}

		cmd.SetContext(context.WithValue(cmd.Context(), constants.ConfigKeyGlobalConfig, GlobalConfig))
		cmd.SetContext(context.WithValue(cmd.Context(), constants.ConfigKeyViper, Viper))
	},
}

func main() {
	initConfig()
	initFlags()

	RegisterCommands()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
