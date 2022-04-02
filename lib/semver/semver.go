package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Semver [3]int

func MustParse(versionString string) (Semver) {
	if matched, _ := regexp.Match(`^(v)?([0-9])+\.([0-9])+\.([0-9])+$`, []byte(versionString)); !matched {
		panic(errors.New("cannot parse version string since it is not in a valid semver format"))
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
	return fmt.Sprintf("%v.x", s[0])
}

func (s Semver) String() string {
	return fmt.Sprintf("%v.%v.%v", s[0], s[1], s[2])
}