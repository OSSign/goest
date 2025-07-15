//go:build wasm || wasi
// +build wasm wasi

package constants

import (
	"github.com/sassoftware/relic/v8/config"
)

type Vault interface {
	GetType() string

	GetTokenConfig() *config.TokenConfig
	GetKeyConfig() *config.KeyConfig
}

type SigningConfig struct {
	EnabledSigners []string `mapstructure:"signers"`
	TimestampUrls  []string `mapstructure:"timestampUrls"`
	AzureKeyVault  Vault    `mapstructure:"azurekv"`
	OutSuffix      string   `mapstructure:"outSuffix"`
}

type NoPassPrompt struct{}

func (NoPassPrompt) GetPasswd(prompt string) (string, error) {
	return "", nil
}
