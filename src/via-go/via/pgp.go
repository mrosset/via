package via

import (
	pgp "code.google.com/p/go.crypto/openpgp"
	"os"
	"os/exec"
	"path/filepath"
)

func Sign(path string) (err error) {
	sig := path + ".sig"
	if fileExists(sig) {
		err = os.Remove(sig)
		if err != nil {
			return err
		}
	}
	return exec.Command("gpg", "-u", "test@test.com",
		"--detach-sign", path).Run()
}

func CheckSig(path string) (err error) {
	keyring := filepath.Join(os.Getenv("HOME"), ".gnupg", "secring.gpg")
	fd, err := os.Open(keyring)
	if err != nil {
		return err
	}
	defer fd.Close()
	keys, err := pgp.ReadKeyRing(fd)
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
	_, err = pgp.CheckDetachedSignature(keys, pkg, sig)
	if err != nil {
		return err
	}
	return nil
}
