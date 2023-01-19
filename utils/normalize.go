package utils

import (
	"errors"
	"regexp"
)

var (
	cellphonePattern *regexp.Regexp
	passwordPattern  *regexp.Regexp
)

func init() {
	cellphonePattern = regexp.MustCompile(`^\d{10}$`)
	passwordPattern = regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
}

func NormalizeCellphone(cellphone string) (string, error) {
	var start int
	if len(cellphone) < 10 || len(cellphone) > 14 {
		return "", errors.New("not a valid cellphone")
	}

	switch {
	case cellphone[:5] == "00989":
		start = 4
	case cellphone[:4] == "+989":
		start = 3
	case cellphone[:3] == "989":
		start = 2
	case cellphone[:2] == "09":
		start = 1
	case cellphone[0] == '9':
		start = 0
	default:
		return "", errors.New("not a valid cellphone")
	}

	if !cellphonePattern.Match([]byte(cellphone[start:])) {
		return "", errors.New("not a valid cellphone")
	}

	cellphone = "0" + cellphone[start:]

	return cellphone, nil
}

func ValidatePassword(password string) error {
	matched := passwordPattern.MatchString(password)
	if !matched {
		return errors.New("password: minimum eight characters of letter and numbers")
	}

	return nil
}
