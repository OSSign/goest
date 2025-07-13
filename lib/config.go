package lib

import "github.com/spf13/viper"

var Viper *viper.Viper

type ConfigKey interface{}

const (
	ConfigKeyFile = "configFile"
)

type SigningConfig struct {
	EnabledSigners []string       `mapstructure:"signers"`
	TimestampUrls  []string       `mapstructure:"timestampUrls"`
	AzureKeyVault  *AzureKeyVault `mapstructure:"azurekv"`
}

var GlobalConfig *SigningConfig

func InitConfig() {
	if Viper == nil {
		Viper = viper.New()
	}

	GlobalConfig = &SigningConfig{
		EnabledSigners: []string{"all"},
		AzureKeyVault:  &AzureKeyVault{},
	}
}
