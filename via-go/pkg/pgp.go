package via

import (
	pgp "crypto/openpgp"
	"exec"
	"os"
	"path/filepath"
)

func Sign(path string) (err os.Error) {
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

func CheckSig(path string) (err os.Error) {
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
