package iban

import (
	"errors"
	"math/big"
	"strconv"
)

var (
	errFailed      = errors.New("IBAN check digits validation failed")
	errCheckDigits = errors.New("IBAN has incorrect check digits")
)

func validateCheckDigits(iban string) error {
	// Move the four initial characters to the end of the string
	iban = iban[4:] + iban[:4]

	// Replace each letter in the string with two digits, thereby expanding the string, where A = 10, B = 11, ..., Z = 35
	mods := ""
	for _, c := range iban {
		// Get character code point value
		i := int(c)

		// Check if c is characters A-Z (codepoint 65 - 90)
		if i > 64 && i < 91 {
			// A=10, B=11 etc...
			i -= 55
			// Add int as string to mod string
			mods += strconv.Itoa(i)
		} else {
			mods += string(c)
		}
	}

	// Create bignum from mod string and perform module
	bigVal, success := new(big.Int).SetString(mods, 10)
	if !success {
		return errFailed
	}

	modVal := new(big.Int).SetInt64(97)
	resVal := new(big.Int).Mod(bigVal, modVal)

	// Check if module is equal to 1
	if resVal.Int64() != 1 {
		return errCheckDigits
	}

	return nil
}
