package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ossign/gost/lib"
	"github.com/ossign/gost/lib/constants"
	"github.com/sassoftware/relic/v8/cmdline/shared"
	"github.com/sassoftware/relic/v8/config"
	"github.com/sassoftware/relic/v8/lib/x509tools"
	"github.com/sassoftware/relic/v8/signers"
	"github.com/sassoftware/relic/v8/token/open"
	"github.com/spf13/cobra"
)

type NoPassPrompt struct{}

func (NoPassPrompt) GetPasswd(prompt string) (string, error) {
	return "", nil // No password prompt, return empty string
}

func RegisterCommands() {
	sign := &cobra.Command{
		Use:   "sign [vault] [args] [file1] [file2] ...",
		Short: "Sign files with the specified vault",
		Long:  "Sign files using the enabled vault specified in the configuration or command line arguments.",
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Signing files: %v\n", args)

			vault, ok := cmd.Context().Value(constants.ConfigKeyVault).(constants.Vault)
			if !ok {
				fmt.Println("No vault specified. ")
				cmd.Help()
				return
			}

			shared.CurrentConfig = &config.Config{
				Tokens: map[string]*config.TokenConfig{
					vault.GetType(): vault.GetTokenConfig(),
				},
				Keys: map[string]*config.KeyConfig{
					vault.GetType(): vault.GetKeyConfig(),
				},
				Server:  nil,
				Clients: nil,
				Remote:  nil,
				Notary:  nil,
				Timestamp: &config.TimestampConfig{
					URLs:   GlobalConfig.TimestampUrls,
					MsURLs: GlobalConfig.TimestampUrls,
				},
			}

			for _, file := range args {
				signr, err := signers.ByFile(file, "")
				if err != nil {
					panic(err)
				}

				if signr.Sign == nil {
					panic("No signing method available for the specified file type")
				}

				hash := x509tools.HashByName("SHA-256")

				token, err := open.Token(shared.CurrentConfig, vault.GetType(), new(NoPassPrompt))

				if err != nil {
					panic(err)
				}

				cert, opts, err := lib.Init(cmd.Context(), signr, token, vault.GetType(), hash, &signers.FlagValues{
					Defs:   &lib.DefaultFlagset,
					Values: map[string]string{},
				})
				if err != nil {
					panic(err)
				}

				opts.Path = file
				infile, err := shared.OpenForPatching(file, file)
				if err != nil {
					panic(err)
				}

				defer infile.Close()

				transform, err := signr.GetTransform(infile, *opts)
				if err != nil {
					panic(err)
				}

				stream, err := transform.GetReader()
				if err != nil {
					panic(err)
				}

				blob, err := signr.Sign(stream, cert, *opts)
				if err != nil {
					panic(err)
				}

				outFilename := file
				if GlobalConfig.OutSuffix != "" {
					splitFilename := strings.Split(file, ".")
					if len(splitFilename) > 1 {
						outFilename = fmt.Sprintf("%s-%s.%s", strings.Join(splitFilename[:len(splitFilename)-1], "."), GlobalConfig.OutSuffix, splitFilename[len(splitFilename)-1])
					} else {
						outFilename = fmt.Sprintf("%s-%s", file, GlobalConfig.OutSuffix)
					}
				}

				mimeType := opts.Audit.GetMimeType()
				if err := transform.Apply(outFilename, mimeType, bytes.NewReader(blob)); err != nil {
					panic(err)
				}

				fmt.Fprintln(os.Stderr, "Signed", file, "and saved to", file)
			}
		},
	}

	sign.AddCommand(GlobalConfig.AzureKeyVault.GetCommand())
	sign.PersistentFlags().StringVar(&GlobalConfig.OutSuffix, "out-suffix", "", "Suffix to append to signed files")

	rootCmd.AddCommand(sign)
}
