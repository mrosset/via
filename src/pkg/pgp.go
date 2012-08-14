package via

import (
	"code.google.com/p/go.crypto/openpgp"
	"code.google.com/p/go.crypto/openpgp/packet"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var keyring = filepath.Join(os.Getenv("HOME"), ".gnupg", "secring.gpg")

func Sign(plan *Plan) (err error) {
	var (
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
		i, ok := k.Identities[config.Identity]
		if ok {
			entity = k
			identity = i
		}
	}
	if entity == nil || identity == nil {
		return fmt.Errorf("Could not find entity or identity for %s", config.Identity)
	}
	if entity.PrivateKey.Encrypted {
		pw := ""
		fmt.Printf("%s Password: ", identity.Name)
		_, err := fmt.Scanln(&pw)
		if err != nil {
			return err
		}
		err = entity.PrivateKey.Decrypt([]byte(pw))
		if err != nil {
			return err
		}
	}
	ppath := path.Join(config.Repo, plan.PackageFile())
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

func CheckSig(path string) (err error) {
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
