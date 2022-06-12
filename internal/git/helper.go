package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"sykesdev.ca/gog/config"
	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/common/constants"
	"sykesdev.ca/gog/internal/logging"
	"sykesdev.ca/gog/internal/semver"
)

func repositoryIsValid() bool {
	cmd := exec.Command("git", "status")
	_, err := cmd.Output()
	
	logging.Instance().Debugf("valid repository: %t", err == nil)

	return err == nil
}

func localBranchExists(branch *Branch) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git branch | egrep %s", branch.Name))
	_, err := cmd.Output()

	logging.Instance().Debugf("local branch exists: %t", err == nil)

	return err == nil
}

func remoteBranchExists(branch *Branch) bool {
	cmd := exec.Command("bash", "-c", "git ls-remote --head origin | egrep " + branch.Name)
	_, err := cmd.Output()

	logging.Instance().Debugf("remote branch exists: %t", err == nil)

	return err == nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("bash", "-c", "git branch | grep '*' | cut -d' ' -f2")
	stdout, err := cmd.CombinedOutput()

	return common.CleanStdoutSingleline(stdout), err
}

func originDefaultBranch() (string, error) {
	defaultBranchCmd := exec.Command("bash", "-c", "git remote show origin | sed -n '/HEAD branch/s/.*: //p'")
	defaultBranch, err := defaultBranchCmd.CombinedOutput()

	return common.CleanStdoutSingleline(defaultBranch), err
}

func projectExistingVersionPrefix() (string, error) {
	tagName, err := originLatestTagName()
	if err != nil {
		logging.Instance().Debugf("error ocurred when reading latest tagName from repo: %v\n%s", err, tagName)

		if strings.Contains(err.Error(), "128") {
			logging.Instance().Debug("existing tag prefix defaulting to global defaults since no existing tags found on remote")
			return config.AppConfig().TagPrefix(), nil
		}

		return "", fmt.Errorf("could not tag information from remote origin. %v", err)
	}

	logging.Instance().Debugf("origin current/latest tagName: %s", tagName)

	var existingPrefix string
	if prefixSearch := regexp.MustCompile(constants.VersionPrefixRegexp).FindStringSubmatch(tagName); len(prefixSearch) > 0 {
		existingPrefix = strings.TrimSpace(prefixSearch[0])
	} else {
		existingPrefix = ""
	}

	logging.Instance().Debugf("captured existing prefix for repository: %s", existingPrefix)

	return existingPrefix, nil
}

func originLatestFullVersion() (semver.Semver, error) {
	version := semver.Semver{0,0,0}

	defaultBranch, err := originDefaultBranch()
	if err != nil {
		return version, err
	}

	logging.Instance().Debugf("default branch at: %s", defaultBranch)

	tagCmd := exec.Command("bash", "-c", fmt.Sprintf("git tag --merged %s", defaultBranch))
	tagOut, err := tagCmd.CombinedOutput()
	if err != nil {
		
		logging.Instance().Debugf("error ocurred when capturing current tag version from remote (%s): %v", defaultBranch, err)

		if strings.Contains(err.Error(), "128") {
			logging.Instance().Debug("defaulting to verion 0.0.0 since no existing tags found on remote")
			return version, nil
		}
		
		return version, err
	}

	semverRegex, err := regexp.Compile(constants.FullSemverRegexp)
	if err != nil {
		return version, err
	}

	logging.Instance().Debug("checking for latest existing tag from remote")

	latestTag := semver.Semver{0,0,0}
	tagScanner := bufio.NewScanner(bytes.NewReader(tagOut))
	for tagScanner.Scan() {
		tag := tagScanner.Text()

		logging.Instance().Debugf("processing: %s", tag)

		if matched := semverRegex.MatchString(tag); matched {
			semverTag, err := semver.Parse(tag)
			if err != nil {
				return version, err
			}

			if semverTag.GreaterThan(latestTag) {
				logging.Instance().Debugf("found newer tag version: %s", semverTag)
				latestTag = semverTag
			}
		}
	}

	logging.Instance().Debugf("latest tag found: %s", latestTag)

	return latestTag, nil
}

func deleteBranch(branch *Branch) error {
	cmdLocal := exec.Command("git", "branch", "-D", branch.Name)
	localStdout, err := cmdLocal.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(localStdout))
	}

	logging.Instance().Debugf("deleted local branch: %s", branch.Name)

	cmdRemote := exec.Command("git", "push", "origin", "--delete", branch.Name)
	remoteStdout, err := cmdRemote.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(remoteStdout))
	}

	logging.Instance().Debugf("deleted remote branch: %s", branch.Name)

	return nil
}

func originLatestTagName() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	stdout, err := cmd.CombinedOutput()

	return common.CleanstdoutMultiline(stdout), err
}