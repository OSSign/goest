//go:build wasm
// +build wasm

package vaults

import (
	"os"

	"github.com/sassoftware/relic/v8/config"
	_ "github.com/sassoftware/relic/v8/token/azuretoken"
)

type AzureKeyVault struct {
	Url    string `mapstructure:"url"`
	Tenant string `mapstructure:"tenant"`
	Client string `mapstructure:"client"`
	Secret string `mapstructure:"secret"`
}

func (v *AzureKeyVault) GetType() string {
	return "azure"
}

func (v *AzureKeyVault) GetTokenConfig() *config.TokenConfig {
	return &config.TokenConfig{
		Type: v.GetType(),
	}
}

func (v *AzureKeyVault) GetKeyConfig() *config.KeyConfig {
	os.Setenv("AZURE_TENANT_ID", v.Tenant)
	os.Setenv("AZURE_CLIENT_ID", v.Client)
	os.Setenv("AZURE_CLIENT_SECRET", v.Secret)

	return &config.KeyConfig{
		Token:     v.GetType(),
		ID:        v.Url,
		Timestamp: true,
	}
}
