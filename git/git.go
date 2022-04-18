package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"sykesdev.ca/gog/common/constants"
	"sykesdev.ca/gog/semver"
)

func HasUncommittedChanges() bool {
	cmd := exec.Command("bash", "-c", "git status --porcelain | egrep '^[A,M,D]'")
	_, err := cmd.Output()

	return err == nil
}

func IsValidRepo() bool {
	cmd := exec.Command("git", "status")
	_, err := cmd.Output()
	
	return err == nil
}

func BranchExists(branch string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git branch | egrep %s", branch))
	_, err := cmd.Output()

	return err == nil
}

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("bash", "-c", "git branch | grep '*' | cut -d' ' -f2")
	stdout, err := cmd.CombinedOutput()

	return string(stdout), err
}

func PullChanges() (string, error) {
	cmd := exec.Command("git", "pull")
	stdout, err := cmd.CombinedOutput()
	
	return string(stdout), err
}

func OriginDefaultBranch() (string, error) {
	defaultBranchCmd := exec.Command("bash", "-c", "git remote show origin | sed -n '/HEAD branch/s/.*: //p'")
	defaultBranch, err := defaultBranchCmd.CombinedOutput()

	return strings.TrimSpace(string(defaultBranch)), err
}

func CheckoutDefaultBranch() (string, error) {
	defaultBranch, err := OriginDefaultBranch()
	if err != nil {
		return defaultBranch, err
	}

	checkoutCmd := exec.Command("git", "checkout", defaultBranch)
	checkoutStdout, err := checkoutCmd.CombinedOutput()
	if err != nil {
		return string(checkoutStdout), err
	}

	if stderr, err := PullChanges(); err != nil {
		return stderr, err
	}

	return "", nil
}

func StageChanges() (string, error) {
	cmd := exec.Command("git", "add", "-A")
	stderr, err := cmd.CombinedOutput()
	
	return string(stderr), err
}

func PushRemote(pushArgs string) (string, error) {
	pushCommand := fmt.Sprintf("git push %s", pushArgs)
	cmd := exec.Command("bash", "-c", pushCommand)
	stderr, err := cmd.CombinedOutput()

	return string(stderr), err
}

func PushRemoteTagsOnly() (string, error) {
	cmd := exec.Command("git", "push", "--tags", "--force")
	stderr, err := cmd.CombinedOutput()
	
	return string(stderr), err
}

func LatestTagName() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func OriginCurrentVersion() (semver.Semver, error) {
	version := semver.Semver{0,0,0}

	defaultBranch, err := OriginDefaultBranch()
	if err != nil {
		return version, err
	}

	tagCmd := exec.Command("bash", "-c", fmt.Sprintf("git tag --merged %s", defaultBranch))
	tagOut, err := tagCmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "128") {
			return version, nil
		}
		
		return version, err
	}

	semverRegex, err := regexp.Compile(constants.FullSemverRegexp)
	if err != nil {
		return version, err
	}

	latestTag := [3]int{0,0,0}
	tagScanner := bufio.NewScanner(bytes.NewReader(tagOut))
	for tagScanner.Scan() {
		tag := tagScanner.Text()
		if matched := semverRegex.MatchString(tag); matched {
			semverTag, err := semver.Parse(tag)
			if err != nil {
				return version, err
			}

			if semverTag.GreaterThan(latestTag) {
				latestTag = semverTag
			}
		}
	}

	return latestTag, nil
}