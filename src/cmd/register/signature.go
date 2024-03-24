package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ProtonMail/go-crypto/openpgp"
)

const (
	sigFileSuffix = ".sig"
)

func checkSig(dir string, el openpgp.EntityList) error {
	p := path.Join(dir, shaFilePattern)
	m, err := filepath.Glob(p)
	if err != nil {
		return fmt.Errorf("failed globbing shasum file with pattern %q - %w", p, err)
	}

	if len(m) != 1 {
		return fmt.Errorf("file globbing pattern %q should have matched exactly 1 file, got %d files", p, len(m))
	}

	payload, err := os.Open(m[0])
	if err != nil {
		return fmt.Errorf("failed to open hash file %q - %w", m[0], err)
	}
	defer func() {
		_ = payload.Close()
	}()

	signature, err := os.Open(m[0] + sigFileSuffix)
	if err != nil {
		return fmt.Errorf("failed to open signature file %q - %w", m[0]+sigFileSuffix, err)
	}
	defer func() {
		_ = signature.Close()
	}()

	signer, err := openpgp.CheckDetachedSignature(el, payload, signature, nil)
	if err != nil {
		return fmt.Errorf("failed signature check - %w", err)
	}

	if len(el.KeysById(signer.PrimaryKey.KeyId)) == 0 {
		return fmt.Errorf("signature okay, but we don't recognize the signer")
	}

	return nil
}
