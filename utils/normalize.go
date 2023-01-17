package utils

import (
	"errors"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	cellphonePattern *regexp.Regexp
	emailPattern     *regexp.Regexp
)

func init() {
	var err error
	emailPattern = regexp.MustCompile(`(^[a-zA-Z0-9.!#$%&'*+=?^_{|}~-]{3,64})@([a-zA-Z0-9.\-]{2,186})\.([a-zA-Z]{2,4})$`)
	cellphonePattern, err = regexp.Compile(`^\d{10}$`)
	if err != nil {
		log.Fatal("cannot compile cellphone regex")
	}
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

func NormalizeEmail(email string) (string, error) {
	matched := emailPattern.MatchString(email)
	if !matched {
		return "", errors.New("not a valid email")
	}

	parts := strings.Split(email, "@")
	parts[0] = strings.ToLower(parts[0])
	parts[1] = strings.ToLower(parts[1])

	return strings.Join(parts, "@"), nil
}
