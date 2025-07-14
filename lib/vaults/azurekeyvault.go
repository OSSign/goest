package vaults

import (
	"context"
	"fmt"
	"os"

	"github.com/ossign/gost/lib/constants"
	"github.com/sassoftware/relic/v8/config"
	_ "github.com/sassoftware/relic/v8/token/azuretoken"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func (v *AzureKeyVault) GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "azurekv [flags] [file1] [file2] ...",
		Short: "Azure Key Vault",
		Long:  "Azure Key Vault is a cloud service for securely storing and accessing secrets, keys, and certificates.",
		Run: func(cmd *cobra.Command, args []string) {
			if v.Url == "" || v.Tenant == "" || v.Client == "" || v.Secret == "" || len(args) == 0 {
				cmd.Help()
				fmt.Println("")
			}

			if v.Url == "" {
				fmt.Println("Azure Key Vault URL is required. Use --url, -u or env GOST_AZUREKV_URL to specify the URL. Example: --url https://myvault.vault.azure.net/certificates/MyCertificate/1234567890abcdef1234567890abcdef")
			}

			if v.Tenant == "" {
				fmt.Println("Azure Tenant ID is required. Use --tenant, -t or env GOST_AZUREKV_TENANT to specify the Tenant ID. Example: --tenant 12345678-1234-1234-1234-123456789012")
			}
			if v.Client == "" {
				fmt.Println("Azure Client ID is required. Use --client, -c or env GOST_AZUREKV_CLIENT to specify the Client ID. Example: --client 12345678-1234-1234-1234-123456789012")
			}
			if v.Secret == "" {
				fmt.Println("Azure Key Vault Secret Name is required. Use --secret, -s or GOST_AZUREKV_SECRET to specify the Secret Name. Example: --secret MySecretName")
			}

			if len(args) == 0 {
				fmt.Println("No files specified for signing. Please provide at least one file to sign.")
				return
			}

			cmd.SetContext(context.WithValue(cmd.Context(), constants.ConfigKeyVault, v))
		},
	}

	cmd.Flags().StringVarP(&v.Url, "url", "u", "", "Azure Key Vault URL (alt. AKV_URL)")
	viper.BindPFlag("akvurl", cmd.Flags().Lookup("url"))

	cmd.Flags().StringVarP(&v.Tenant, "tenant", "t", "", "Azure Tenant ID (alt. AKV_TENANT)")
	viper.BindPFlag("akvTenant", cmd.Flags().Lookup("tenant"))

	cmd.Flags().StringVarP(&v.Client, "client", "c", "", "Azure Client ID (alt. AKV_CLIENT)")
	viper.BindPFlag("akvClient", cmd.Flags().Lookup("client"))

	cmd.Flags().StringVarP(&v.Secret, "secret", "s", "", "Azure Key Vault Secret Name (alt. AKV_SECRET)")
	viper.BindPFlag("akvSecret", cmd.Flags().Lookup("secret"))

	return cmd
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
