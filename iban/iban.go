package iban

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// IBAN struct
type IBAN struct {
	// Full code
	Code string

	// Full code prettyfied for printing on paper
	PrintCode string

	// Country code
	CountryCode string

	// Check digits
	CheckDigits string

	// Country specific bban part
	BBAN string
}

// need to handle all kinds of spaces: [\t\n\f\r ]
var (
	allWhiteSpaces    = regexp.MustCompile(`\s+`)
	numUpperOnly      = regexp.MustCompile(`^[0-9A-Z]*$`)
	twoUpperTwoDigits = regexp.MustCompile(`^[[:upper:]]{2}\d{2}`)

	errAlphaNumOnly             = errors.New("IBAN can contain only alphanumeric characters")
	errCountryCodeOrCheckDigits = errors.New("IBAN must start with country code (2 characters) and check digits (2 digits)")

	// there might be new countries join the IBAN standard, so when got this error, better to double confirm with another IBAN validation service (like iban.com)
	ErrUnsupportedCountryCode = errors.New("Unsupported country code")
)

// NewIBAN create new IBAN with validation
func ParseIBAN(s string) (*IBAN, error) {
	// Prepare string: remove spaces and convert to upper case
	s = allWhiteSpaces.ReplaceAllString(s, "")
	s = strings.ToUpper(s)
	code := s

	// Validate characters
	if !numUpperOnly.MatchString(s) {
		return nil, errAlphaNumOnly
	}

	// Get country code and check digits
	hs := twoUpperTwoDigits.FindString(s)
	if hs == "" {
		return nil, errCountryCodeOrCheckDigits
	}

	countryCode := hs[0:2]
	checkDigits := hs[2:4]

	// Get country settings for country code
	cs, ok := countries[countryCode]
	if !ok {
		return nil, ErrUnsupportedCountryCode
		// fmt.Errorf("Unsupported country code %v", iban.CountryCode)
	}

	// Validate code length
	if len(s) != cs.Length {
		return nil, fmt.Errorf("IBAN length %d does not match length %d specified for country code %v", len(s), cs.Length, countryCode)
	}

	// Set and validate BBAN part, the part after the language code and check digits
	bban := s[4:]

	if err := cs.validateBasicBankAccountNumber(bban); err != nil {
		return nil, err
	}

	// Validate check digits with mod97
	if err := validateCheckDigits(code); err != nil {
		return nil, err
	}

	// Generate print code from code (splits code in sections of 4 characters)
	prc := ""
	for len(s) > 4 {
		prc += s[:4] + " "
		s = s[4:]
	}
	prc += s

	return &IBAN{
		Code:        code,
		PrintCode:   prc,
		CountryCode: countryCode,
		CheckDigits: checkDigits,
		BBAN:        bban,
	}, nil
}
