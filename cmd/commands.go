package main

import (
	"github.com/ossign/gost/lib"
	"github.com/spf13/cobra"
)

func RegisterCommands() {
	sign := &cobra.Command{
		Use:   "sign [vault] [args] [file1] [file2] ...",
		Short: "Sign files with the specified signers",
		Long:  "Sign files using the enabled signers specified in the configuration or command line arguments.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	sign.AddCommand(lib.GlobalConfig.AzureKeyVault.GetCommand())

	rootCmd.AddCommand(sign)
}
