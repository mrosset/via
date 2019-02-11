package via

import (
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"os"
	"path"
	"path/filepath"
)

//TODO: use our own keyring
var keyring = filepath.Join(os.Getenv("HOME"), ".gnupg", "secring.gpg")

// Sign produces a detached signature for a plans package file
//
// FIXME: currently this is not used at all. And may be redundant
// considering we are using ipfs multihash. Will revist this later.
func Sign(ctx *PlanContext) (err error) {
	var (
		plan     = ctx.Plan
		entity   *openpgp.Entity
		identity *openpgp.Identity
	)
	fd, err := os.Open(keyring)
	if err != nil {
		return err
	}
	keys, err := openpgp.ReadKeyRing(fd)
	if err != nil {
		return err
	}
	for _, k := range keys {
		i, ok := k.Identities[ctx.Config.Identity]
		if ok {
			entity = k
			identity = i
		}
	}
	if entity == nil || identity == nil {
		return fmt.Errorf("Could not find entity or identity for %s", ctx.Config.Identity)
	}
	if entity.PrivateKey.Encrypted {
		// TODO: prompt for user Password use keyagent?
		pw := "test"
		/*
			fmt.Printf("%s Password: ", identity.Name)
			_, err := fmt.Scanln(&pw)
			if err != nil {
				return err
			}
		*/
		err = entity.PrivateKey.Decrypt([]byte(pw))
		if err != nil {
			return err
		}
	}
	ppath := path.Join(ctx.Config.Repo, plan.PackageFile())
	pkg, err := os.Open(ppath)
	if err != nil {
		return err
	}
	defer pkg.Close()
	sig, err := os.Create(ppath + ".sig")
	if err != nil {
		return err
	}
	defer sig.Close()
	fmt.Printf(lfmt, "signing", plan.PackageFile())
	err = openpgp.DetachSign(sig, entity, pkg, new(packet.Config))
	if err != nil {
		return err
	}
	return nil
}

// VerifiySig verifiy's that the plans signature matches a trusted
// signature.
//
// FIXME: currently this is not used right now. instead ipfs's
// multihash is used. Revisit this later
func VerifiySig(path string) (err error) {
	fd, err := os.Open(keyring)
	if err != nil {
		return err
	}
	defer fd.Close()
	keys, err := openpgp.ReadKeyRing(fd)
	if err != nil {
		return err
	}
	pkg, err := os.Open(path)
	if err != nil {
		return err
	}
	sig, err := os.Open(path + ".sig")
	if err != nil {
		return err
	}
	_, err = openpgp.CheckDetachedSignature(keys, pkg, sig)
	if err != nil {
		return err
	}
	return nil
}
