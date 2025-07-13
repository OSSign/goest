package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ossign/gost/lib"
	"github.com/spf13/cobra"
)

var signers []string = []string{
	"appx", "cab", "cat", "deb", "jar", "msi", "pecoff", "pgp", "pkcs", "ps", "rpm", "vsix", "xap",
}

var timestampUrls []string = []string{
	"http://timestamp.globalsign.com/tsa/advanced",
}

var configFile string

var rootCmd = &cobra.Command{
	Use:   "gost",
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

			lib.Viper.SetConfigFile(configFile)
			lib.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			lib.Viper.SetEnvPrefix("gost")
			lib.Viper.AutomaticEnv()

			if err := lib.Viper.ReadInConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Using config file: %s\n", lib.Viper.ConfigFileUsed())

			if err := lib.Viper.Unmarshal(&lib.GlobalConfig); err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshalling config: %v\n", err)
				os.Exit(1)
			}

			lib.RegisterSigners(signers)
		}
	},
}

func initFlags() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "f", "", "Path to a configuration file (optional)")
	rootCmd.PersistentFlags().StringSliceVarP(&signers, "signers", "", signers, "List of signers to enable (default: all available signers)")
	rootCmd.PersistentFlags().StringSliceVar(&timestampUrls, "timestampUrls", timestampUrls, "List of timestamp URLs to use for signing (optional, default http://timestamp.globalsign.com/tsa/advanced)")

	lib.Viper.BindPFlag("signers", rootCmd.PersistentFlags().Lookup(("signers")))
	lib.Viper.BindPFlag("timestampUrls", rootCmd.PersistentFlags().Lookup("timestampUrls"))
}

func main() {
	lib.InitConfig()
	initFlags()

	RegisterCommands()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
