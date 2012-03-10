package via

import (
	"code.google.com/p/go.crypto/openpgp"
	"code.google.com/p/go.crypto/openpgp/packet"
	"exp/terminal"
	"fmt"
	"os"
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
	if entity.PrivateKey.Encrypted {
		_ = identity.Name
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
	ppath := filepath.Join(config.Repo, plan.PackageFile())
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
	fmt.Println("Signing", ppath)
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

func ReadPassword() (string, error) {
	fd, err := os.OpenFile("/dev/ttys001", os.O_RDWR, 0)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	term := terminal.NewTerminal(fd, "$")
	line, err := term.ReadPassword("Password: ")
	if err != nil {
		return "", err
	}
	return line, nil
}
