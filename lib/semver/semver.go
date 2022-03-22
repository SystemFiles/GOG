package semver

import (
	"fmt"
	"strconv"
	"strings"
)

type UpdateLevel string
type Semver [3]int

func Parse(versionString string) (Semver) {
	elements := strings.Split(string(versionString), ".")
	major, err := strconv.Atoi(elements[0])
	if err != nil { major = 0 }
	minor, err := strconv.Atoi(elements[1])
	if err != nil { minor = 0 }
	patch, err := strconv.Atoi(elements[2])
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