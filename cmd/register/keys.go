package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/ProtonMail/go-crypto/openpgp"
)

const (
	keyFile = "gpg_key.asc"
)

func getKeyRing(b []byte) (openpgp.EntityList, error) {
	el, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed parsing key armor file contents - %w", err)
	}

	if len(el) != 1 {
		return nil, fmt.Errorf("expected keyring to have 1 element, got %d elements", len(el))
	}

	return el, nil
}

func shortId(el openpgp.EntityList) (string, error) {
	if len(el) != 1 {
		return "", fmt.Errorf("expected keyring to have 1 element, got %d elements", len(el))
	}

	return fmt.Sprintf("%X", el[0].PrimaryKey.Fingerprint[12:20]), nil
}

func getKeys(dir string) (json.RawMessage, error) {
	fileName := path.Join(dir, keyFile)

	asciiArmor, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("failed reading key armor file %q - %s", fileName, err)
	}

	el, err := getKeyRing(asciiArmor)
	if err != nil {
		log.Fatalf("failed parsing key armor data %q - %s", fileName, err)
	}

	keyId, err := shortId(el)
	if err != nil {
		log.Fatalf("failed finding key ID - %s", err)
	}

	err = checkSig(dir, el)
	if err != nil {
		log.Fatal(err)
	}

	type resultKey struct {
		KeyId      string `json:"key_id"`
		AsciiArmor string `json:"ascii_armor"`
	}

	type result struct {
		GpgPublicKeys []resultKey `json:"gpg_public_keys"`
	}

	return json.Marshal(&result{
		GpgPublicKeys: []resultKey{
			{
				KeyId:      keyId,
				AsciiArmor: string(asciiArmor),
			},
		},
	})
}
