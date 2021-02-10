package utils

import (
	"encoding/hex"
	"strings"
)

// IsAddressValid verifies that a given string is a correct hex code for a address
func IsAddressValid(address string) bool {
	sanitizedAddress := strings.Replace(address, "0x", "", 1)
	_, err := hex.DecodeString(sanitizedAddress)
	return err == nil
}
