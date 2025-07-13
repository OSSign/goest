package lib

import (
	"log"

	"github.com/sassoftware/relic/v8/signers"
	"github.com/spf13/pflag"

	"github.com/sassoftware/relic/v8/signers/apk"
	"github.com/sassoftware/relic/v8/signers/appmanifest"
	"github.com/sassoftware/relic/v8/signers/appx"
	"github.com/sassoftware/relic/v8/signers/cab"
	"github.com/sassoftware/relic/v8/signers/cat"
	"github.com/sassoftware/relic/v8/signers/deb"
	"github.com/sassoftware/relic/v8/signers/jar"
	"github.com/sassoftware/relic/v8/signers/msi"
	"github.com/sassoftware/relic/v8/signers/pecoff"
	"github.com/sassoftware/relic/v8/signers/pgp"
	"github.com/sassoftware/relic/v8/signers/pkcs"
	"github.com/sassoftware/relic/v8/signers/ps"
	"github.com/sassoftware/relic/v8/signers/rpm"
	"github.com/sassoftware/relic/v8/signers/vsix"
	"github.com/sassoftware/relic/v8/signers/xap"
)

var registerSigners = map[string]*signers.Signer{
	"apk":         apk.ApkSigner,
	"appmanifest": appmanifest.AppSigner,
	"appx":        appx.AppxSigner,
	"cab":         cab.CabSigner,
	"cat":         cat.CatSigner,
	"deb":         deb.DebSigner,
	"jar":         jar.JarSigner,
	"msi":         msi.MsiSigner,
	"pecoff":      pecoff.PeSigner,
	"pgp":         pgp.PgpSigner,
	"pkcs":        pkcs.PkcsSigner,
	"ps":          ps.PsSigner,
	"rpm":         rpm.RpmSigner,
	"vsix":        vsix.Signer,
	"xap":         xap.XapSigner,
}

var DefaultFlagset = pflag.FlagSet{}
var RegisteredSigners = []string{}

func RegisterAllSigners() {
	for name, signer := range registerSigners {
		signers.Register(signer)
		DefaultFlagset.AddFlagSet(signer.Flags())
		RegisteredSigners = append(RegisteredSigners, name)
	}
}

func RegisterSigner(name string) {
	if signer, ok := registerSigners[name]; ok {
		signers.Register(signer)
		DefaultFlagset.AddFlagSet(signer.Flags())
		RegisteredSigners = append(RegisteredSigners, name)
	} else {
		log.Fatalf("Unknown signer: %s", name)
	}
}

func RegisterSigners(names []string) {
	for _, name := range names {
		RegisterSigner(name)
	}
}
