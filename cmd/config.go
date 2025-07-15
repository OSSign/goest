//go:build !wasm

package main

import (
	"github.com/ossign/goest/lib/constants"
	"github.com/ossign/goest/lib/vaults"
	"github.com/spf13/viper"
)

var Viper *viper.Viper
var GlobalConfig *constants.SigningConfig

var defaultSigners []string = []string{
	"appx", "cab", "cat", "deb", "jar", "msi", "pecoff", "pgp", "pkcs", "ps", "rpm", "vsix", "xap",
}

var defaultTimestampUrls []string = []string{
	"http://timestamp.globalsign.com/tsa/advanced",
}

func initConfig() {
	if Viper == nil {
		Viper = viper.New()
	}

	GlobalConfig = &constants.SigningConfig{
		EnabledSigners: []string{"all"},
		AzureKeyVault:  &vaults.AzureKeyVault{},
	}
}

func initFlags() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "f", "", "Path to a configuration file (optional)")
	rootCmd.PersistentFlags().StringSliceVarP(&defaultSigners, "signers", "", defaultSigners, "List of signers to enable (default: all available signers)")
	rootCmd.PersistentFlags().StringSliceVar(&defaultTimestampUrls, "timestampUrls", defaultTimestampUrls, "List of timestamp URLs to use for signing (optional, default http://timestamp.globalsign.com/tsa/advanced)")

	Viper.BindPFlag("signers", rootCmd.PersistentFlags().Lookup(("signers")))
	Viper.BindPFlag("timestampUrls", rootCmd.PersistentFlags().Lookup("timestampUrls"))
}
