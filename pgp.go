package via

import (
	"fmt"
	"openpgp"
	"os"
	"path/filepath"
	"util"
)

var keyring = filepath.Join(os.Getenv("HOME"), ".gnupg", "secring.gpg")

func Sign(paths []string) (err error) {
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
	if entity.PrivateKey.Encrypted {
		pw := ""
		fmt.Printf("%s Password: ", identity.Name)
		_, err := fmt.Scanln(&pw)
		util.CheckFatal(err)
		err = entity.PrivateKey.Decrypt([]byte(pw))
		util.CheckFatal(err)
	}
	for _, path := range paths {
		pkg, err := os.Open(path)
		if err != nil {
			return err
		}
		defer pkg.Close()
		sig, err := os.Create(path + ".sig")
		if err != nil {
			return err
		}
		defer sig.Close()
		fmt.Println("Signing", path)
		err = openpgp.DetachSign(sig, entity, pkg)
		if err != nil {
			return err
		}
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
