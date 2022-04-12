package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"sykesdev.ca/gog/config"
)

type Semver [3]int

func isValidSemver(versionString string) bool {
	matched, _ := regexp.Match(`^(v)?([0-9])+\.([0-9])+\.([0-9])+$`, []byte(versionString))
	return matched
}

func Parse(versionString string) (Semver, error) {
	if !isValidSemver(versionString) {
		return Semver{}, errors.New("cannot parse version string provided since it is not in a valid semver format")
	}

	return MustParse(versionString), nil
}

func MustParse(versionString string) (Semver) {
	if !isValidSemver(versionString) {
		panic("cannot parse version string since it is not in a valid semver format")
	}

	numReg := regexp.MustCompile(`[0-9]+`)
	elements := strings.Split(string(versionString), ".")
	major, err := strconv.Atoi(numReg.FindString(elements[0]))
	if err != nil { major = 0 }
	minor, err := strconv.Atoi(numReg.FindString(elements[1]))
	if err != nil { minor = 0 }
	patch, err := strconv.Atoi(numReg.FindString(elements[2]))
	if err != nil { patch = 0 }

	return [3]int{major, minor, patch}
}

func (s Semver) BumpMajor() Semver {
	s[0] += 1
	s[1], s[2] = 0, 0
	return s
}

func (s Semver) BumpMinor() Semver {
	s[1] += 1
	s[2] = 0
	return s
}

func (s Semver) BumpPatch() Semver {
	s[2] += 1
	return s
}

func (s Semver) Major() string {
	return fmt.Sprintf("%s%v.x", config.AppConfig().TagPrefix(), s[0])
}

func (s Semver) String() string {
	return fmt.Sprintf("%s%v.%v.%v", config.AppConfig().TagPrefix(), s[0], s[1], s[2])
}

func (s Semver) NoPrefix() string {
	return fmt.Sprintf("%v.%v.%v", s[0], s[1], s[2])
}

func (s Semver) Equal(o Semver) bool {
	return s[0] == o[0] && s[1] == o[1] && s[2] == o[2]
}

func (s Semver) GreaterThan(o Semver) bool {
	if s[0] > o[0] {
		return true
	}
	
	if s[0] == o[0] && s[1] > o[1] {
		return true
	}

	if s[0] == o[0] && s[1] == o[1] && s[2] > o[2] {
		return true
	}

	return false
}
