package iban

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// CountrySettings contains length for IBAN and format for BBAN
type countrySettings struct {
	// Length of IBAN code for this country
	Length int

	// Format of BBAN part of IBAN for this country
	Format string

	// to save the pre-compiled regexp
	re *regexp.Regexp
}

/*
	Taken from http://www.tbg5-finance.org/ code example
*/
var countries = map[string]*countrySettings{
	"AD": &countrySettings{Length: 24, Format: "F04F04A12"},
	"AE": &countrySettings{Length: 23, Format: "F03F16"},
	"AL": &countrySettings{Length: 28, Format: "F08A16"},
	"AT": &countrySettings{Length: 20, Format: "F05F11"},
	"AZ": &countrySettings{Length: 28, Format: "U04A20"},
	"BA": &countrySettings{Length: 20, Format: "F03F03F08F02"},
	"BE": &countrySettings{Length: 16, Format: "F03F07F02"},
	"BG": &countrySettings{Length: 22, Format: "U04F04F02A08"},
	"BH": &countrySettings{Length: 22, Format: "U04A14"},
	"BR": &countrySettings{Length: 29, Format: "F08F05F10U01A01"},
	"CH": &countrySettings{Length: 21, Format: "F05A12"},
	"CR": &countrySettings{Length: 21, Format: "F03F14"},
	"CY": &countrySettings{Length: 28, Format: "F03F05A16"},
	"CZ": &countrySettings{Length: 24, Format: "F04F06F10"},
	"DE": &countrySettings{Length: 22, Format: "F08F10"},
	"DK": &countrySettings{Length: 18, Format: "F04F09F01"},
	"DO": &countrySettings{Length: 28, Format: "U04F20"},
	"EE": &countrySettings{Length: 20, Format: "F02F02F11F01"},
	"ES": &countrySettings{Length: 24, Format: "F04F04F01F01F10"},
	"FI": &countrySettings{Length: 18, Format: "F06F07F01"},
	"FO": &countrySettings{Length: 18, Format: "F04F09F01"},
	"FR": &countrySettings{Length: 27, Format: "F05F05A11F02"},
	"GB": &countrySettings{Length: 22, Format: "U04F06F08"},
	"GE": &countrySettings{Length: 22, Format: "U02F16"},
	"GI": &countrySettings{Length: 23, Format: "U04A15"},
	"GL": &countrySettings{Length: 18, Format: "F04F09F01"},
	"GR": &countrySettings{Length: 27, Format: "F03F04A16"},
	"GT": &countrySettings{Length: 28, Format: "A04A20"},
	"HR": &countrySettings{Length: 21, Format: "F07F10"},
	"HU": &countrySettings{Length: 28, Format: "F03F04F01F15F01"},
	"IE": &countrySettings{Length: 22, Format: "U04F06F08"},
	"IL": &countrySettings{Length: 23, Format: "F03F03F13"},
	"IS": &countrySettings{Length: 26, Format: "F04F02F06F10"},
	"IT": &countrySettings{Length: 27, Format: "U01F05F05A12"},
	"JO": &countrySettings{Length: 30, Format: "U04F04A18"},
	"KW": &countrySettings{Length: 30, Format: "U04A22"},
	"KZ": &countrySettings{Length: 20, Format: "F03A13"},
	"LB": &countrySettings{Length: 28, Format: "F04A20"},
	"LC": &countrySettings{Length: 32, Format: "U04A24"},
	"LI": &countrySettings{Length: 21, Format: "F05A12"},
	"LT": &countrySettings{Length: 20, Format: "F05F11"},
	"LU": &countrySettings{Length: 20, Format: "F03A13"},
	"LV": &countrySettings{Length: 21, Format: "U04A13"},
	"MC": &countrySettings{Length: 27, Format: "F05F05A11F02"},
	"MD": &countrySettings{Length: 24, Format: "A20"},
	"ME": &countrySettings{Length: 22, Format: "F03F13F02"},
	"MK": &countrySettings{Length: 19, Format: "F03A10F02"},
	"MR": &countrySettings{Length: 27, Format: "F05F05F11F02"},
	"MT": &countrySettings{Length: 31, Format: "U04F05A18"},
	"MU": &countrySettings{Length: 30, Format: "U04F02F02F12F03U03"},
	"NL": &countrySettings{Length: 18, Format: "U04F10"},
	"NO": &countrySettings{Length: 15, Format: "F04F06F01"},
	"PK": &countrySettings{Length: 24, Format: "U04A16"},
	"PL": &countrySettings{Length: 28, Format: "F08F16"},
	"PS": &countrySettings{Length: 29, Format: "U04A21"},
	"PT": &countrySettings{Length: 25, Format: "F04F04F11F02"},
	"QA": &countrySettings{Length: 29, Format: "U04A21"},
	"RO": &countrySettings{Length: 24, Format: "U04A16"},
	"RS": &countrySettings{Length: 22, Format: "F03F13F02"},
	"SA": &countrySettings{Length: 24, Format: "F02A18"},
	"SC": &countrySettings{Length: 31, Format: "U04F02F02F16U03"},
	"SE": &countrySettings{Length: 24, Format: "F03F16F01"},
	"SI": &countrySettings{Length: 19, Format: "F05F08F02"},
	"SK": &countrySettings{Length: 24, Format: "F04F06F10"},
	"SM": &countrySettings{Length: 27, Format: "U01F05F05A12"},
	"ST": &countrySettings{Length: 25, Format: "F08F11F02"},
	"TL": &countrySettings{Length: 23, Format: "F03F14F02"},
	"TN": &countrySettings{Length: 24, Format: "F02F03F13F02"},
	"TR": &countrySettings{Length: 26, Format: "F05A01A16"},
	"UA": &countrySettings{Length: 29, Format: "F06A19"},
	"VG": &countrySettings{Length: 24, Format: "U04F16"},
	"XK": &countrySettings{Length: 20, Format: "F04F10F02"},
}

var errBBANMalformat = errors.New("BBAN part of IBAN is not formatted according to country specification")

func (cs *countrySettings) validateBasicBankAccountNumber(bban string) (err error) {
	if cs.re == nil {
		// the initial one time only
		// Compile the regex and save for later re-use
		if cs.re, err = compileFormat(cs.Format); err != nil {
			return fmt.Errorf("Failed to validate BBAN: %v", err.Error())
		}
	}

	if !cs.re.MatchString(bban) {
		return errBBANMalformat
	}

	return nil
}

var formatParts = regexp.MustCompile(`[ABCFLUW]\d{2}`)

func compileFormat(format string) (*regexp.Regexp, error) {
	// Get format part strings
	fps := formatParts.FindAllString(format, -1)

	// Create regex from format parts
	bbr := ""

	for _, ps := range fps {
		switch ps[:1] {
		case "F":
			bbr += "[0-9]"
		case "L":
			bbr += "[a-z]"
		case "U":
			bbr += "[A-Z]"
		case "A":
			bbr += "[0-9A-Za-z]"
		case "B":
			bbr += "[0-9A-Z]"
		case "C":
			bbr += "[A-Za-z]"
		case "W":
			bbr += "[0-9a-z]"
		}

		// Get repeat factor for group
		repeat, atoiErr := strconv.Atoi(ps[1:])
		if atoiErr != nil {
			return nil, fmt.Errorf("Failed to validate BBAN: %v", atoiErr.Error())
		}

		// Add to regex
		bbr += fmt.Sprintf("{%d}", repeat)
	}

	return regexp.Compile(bbr)
}
