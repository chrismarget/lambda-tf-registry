package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const shaFilePattern = "*_SHA256SUMS"

func getHashFile(dir string) (string, error) {
	p := path.Join(dir, shaFilePattern)
	m, err := filepath.Glob(p)
	if err != nil {
		return "", fmt.Errorf("failed globbing shasum file with pattern %q - %w", p, err)
	}
	if len(m) != 1 {
		return "", fmt.Errorf("file globbing pattern %q should have matched exactly 1 file, got %d files", p, len(m))
	}

	return m[0], nil
}

func hashes(hashFile string) (map[string]string, error) {
	f, err := os.Open(hashFile)
	if err != nil {
		return nil, fmt.Errorf("failed reading hash file %q - %w", hashFile, err)
	}
	defer func() {
		_ = f.Close()
	}()

	result := make(map[string]string)

	var i int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "  ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("checksum file %q line %d unexpected format", hashFile, i+1)
		}

		filename := path.Join(path.Dir(hashFile), parts[1])
		expectedHash := parts[0]

		err = checkHash(filename, expectedHash)
		if err != nil {
			return nil, err
		}

		result[parts[0]] = path.Join(path.Dir(hashFile), parts[1])
		i++
	}

	return result, nil
}

func checkHash(filename, expected string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q - %w", filename, err)
	}
	defer func() {
		_ = f.Close()
	}()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return fmt.Errorf("failed copying from %q to the hash writer - %w", filename, err)
	}

	sha256sum := fmt.Sprintf("%x", h.Sum(nil))
	if expected != sha256sum {
		return fmt.Errorf("file %q expected sha256: %s got %s", filename, expected, sha256sum)
	}

	return nil
}
