package mandosjsonparse

import "golang.org/x/crypto/sha3"

// Keccak256 cryptographic function
// TODO: externalize the same way as the file resolver
func keccak256(data []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	result := hash.Sum(nil)
	return result, nil
}
