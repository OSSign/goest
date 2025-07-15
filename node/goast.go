//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ossign/goest/lib/constants"
	"github.com/ossign/goest/lib/vaults"
	"github.com/sassoftware/relic/v8/cmdline/shared"
	"github.com/sassoftware/relic/v8/config"
	"github.com/sassoftware/relic/v8/signers"
)

// var Running = make(chan int, 1)
var Buffer [12][1024]byte
var BufI int = 0

var GlobalConfig *constants.SigningConfig

func main() {
	// for {
	// 	code := <-Running

	// 	fmt.Println("Received signal to exit the function: ", code)
	// 	os.Exit(code)
	// }
}

//go:wasmexport MkBuffer
func MkBuffer() *byte {
	BufI++
	if BufI >= len(Buffer) {
		panic("Buffer overflow, resetting index. Please increase buffer size to be able to pass more functions.")
	}
	return &Buffer[BufI-1][0]
}

//go:wasmexport Initialize
func Initialize() {
	fmt.Println("Initializing Goest...")
	GlobalConfig = &constants.SigningConfig{
		EnabledSigners: []string{"all"},
		AzureKeyVault:  &vaults.AzureKeyVault{},
	}
}

//go:wasmexport Configure
func Configure(url, tenant, clientID, clientSecret string) {
	if url == "" || tenant == "" || clientID == "" || clientSecret == "" {
		panic("Azure Key Vault configuration is incomplete. Please provide url, tenant, clientID, and clientSecret.")
	}

	fmt.Println("GlobalConfig", GlobalConfig)

	GlobalConfig.EnabledSigners = []string{"appx", "cab", "cat", "deb", "jar", "msi", "pecoff", "pgp", "pkcs", "ps", "vsix", "xap"}
	GlobalConfig.TimestampUrls = []string{"http://timestamp.globalsign.com/tsa/advanced"}
	GlobalConfig.AzureKeyVault = &vaults.AzureKeyVault{
		Tenant: tenant,
		Client: clientID,
		Secret: clientSecret,
		Url:    url,
	}

	fmt.Println("Successfully configured Azure Key Vault")
}

//go:wasmexport Sign
func Sign(files string) {
	fmt.Println("Starting signing...")
	fmt.Println("Listing files...", files)

	// filex, err := os.ReadDir("/data")
	// if err != nil {
	// 	panic(fmt.Sprintf("Error reading directory: %v", err))
	// }

	// fmt.Println("Files to sign:", files, filex)

	content, err := os.ReadFile(files)
	if err != nil {
		fmt.Println("XXX")
		panic(fmt.Sprintf("Error reading file %s: %v", files, err))
	}
	fmt.Println("Content of file:", string(content))

	shared.CurrentConfig = &config.Config{
		Tokens: map[string]*config.TokenConfig{
			GlobalConfig.AzureKeyVault.GetType(): GlobalConfig.AzureKeyVault.GetTokenConfig(),
		},
		Keys: map[string]*config.KeyConfig{
			GlobalConfig.AzureKeyVault.GetType(): GlobalConfig.AzureKeyVault.GetKeyConfig(),
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

	for file := range strings.SplitSeq(files, ";") {
		fmt.Println("Processing file:", file)
		signr, err := signers.ByFile(file, "")
		if err != nil {
			fmt.Println("Error finding signer for file:", file, "Error:")
			os.Exit(1)
		}

		fmt.Println("Found signer:", signr.Name)

		// if signr.Sign == nil {
		// 	panic("No signing method available for the specified file type")
		// }

		// hash := x509tools.HashByName("SHA-256")

		// token, err := open.Token(shared.CurrentConfig, GlobalConfig.AzureKeyVault.GetType(), new(constants.NoPassPrompt))
		// if err != nil {
		// 	panic(err)
		// }

		// cert, opts, err := lib.Init(context.Background(), signr, token, GlobalConfig.AzureKeyVault.GetType(), hash, &signers.FlagValues{
		// 	Defs:   &lib.DefaultFlagset,
		// 	Values: map[string]string{},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// opts.Path = file
		// infile, err := shared.OpenForPatching(file, file)
		// if err != nil {
		// 	panic(err)
		// }

		// defer infile.Close()

		// transform, err := signr.GetTransform(infile, *opts)
		// if err != nil {
		// 	panic(err)
		// }

		// stream, err := transform.GetReader()
		// if err != nil {
		// 	panic(err)
		// }

		// blob, err := signr.Sign(stream, cert, *opts)
		// if err != nil {
		// 	panic(err)
		// }

		// outFilename := file
		// if GlobalConfig.OutSuffix != "" {
		// 	splitFilename := strings.Split(file, ".")
		// 	if len(splitFilename) > 1 {
		// 		outFilename = fmt.Sprintf("%s-%s.%s", strings.Join(splitFilename[:len(splitFilename)-1], "."), GlobalConfig.OutSuffix, splitFilename[len(splitFilename)-1])
		// 	} else {
		// 		outFilename = fmt.Sprintf("%s-%s", file, GlobalConfig.OutSuffix)
		// 	}
		// }

		// mimeType := opts.Audit.GetMimeType()
		// if err := transform.Apply(outFilename, mimeType, bytes.NewReader(blob)); err != nil {
		// 	panic(err)
		// }

		// fmt.Fprintln(os.Stderr, "Signed", file, "and saved to", file)
	}

	// Running <- 0
	// fmt.Println("Done!")
}
