package constants

import (
	"github.com/sassoftware/relic/v8/config"
	"github.com/spf13/cobra"
)

type Vault interface {
	GetType() string
	GetCommand() *cobra.Command

	GetTokenConfig() *config.TokenConfig
	GetKeyConfig() *config.KeyConfig
}

type SigningConfig struct {
	EnabledSigners []string `mapstructure:"signers"`
	TimestampUrls  []string `mapstructure:"timestampUrls"`
	AzureKeyVault  Vault    `mapstructure:"azurekv"`
	OutSuffix      string   `mapstructure:"outSuffix"`
}
