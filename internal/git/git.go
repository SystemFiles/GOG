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

func HasUncommittedChanges() bool {
	cmd := exec.Command("bash", "-c", "git status --porcelain | egrep '^[A,M,D,R]'")
	_, err := cmd.Output()

	logging.Instance().Debugf("uncommitted changes: %b", err == nil)

	return err == nil
}

func IsValidRepo() bool {
	cmd := exec.Command("git", "status")
	_, err := cmd.Output()
	
	logging.Instance().Debugf("valid repository: %b", err == nil)

	return err == nil
}

func LocalBranchExists(branch string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git branch | egrep %s", branch))
	_, err := cmd.Output()

	logging.Instance().Debugf("local branch exists: %b", err == nil)

	return err == nil
}

func RemoteBranchExists(branch string) bool {
	cmd := exec.Command("bash", "-c", "git ls-remote --head origin | egrep " + branch)
	_, err := cmd.Output()

	logging.Instance().Debugf("remote branch exists: %b", err == nil)

	return err == nil
}

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("bash", "-c", "git branch | grep '*' | cut -d' ' -f2")
	stdout, err := cmd.CombinedOutput()

	return common.CleanStdoutSingleline(stdout), err
}

func Commit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	stderr, err := cmd.CombinedOutput()

	return common.CleanstdoutMultiline(stderr), err
}

func PullChanges() (string, error) {
	cmd := exec.Command("git", "pull")
	stdout, err := cmd.CombinedOutput()
	
	return common.CleanstdoutMultiline(stdout), err
}

func OriginDefaultBranch() (string, error) {
	defaultBranchCmd := exec.Command("bash", "-c", "git remote show origin | sed -n '/HEAD branch/s/.*: //p'")
	defaultBranch, err := defaultBranchCmd.CombinedOutput()

	return common.CleanStdoutSingleline(defaultBranch), err
}

func Checkout(branch string, create bool) (string, error) {
	checkoutArgs := make([]string, 0)
	if create {
		checkoutArgs = append(checkoutArgs, "-b")
	}
	checkoutArgs = append(checkoutArgs, branch)

	logging.Instance().Debugf("checking out branch, %s, with create: %b", branch, create)

	cmd := exec.Command("git", append([]string{"checkout"}, checkoutArgs...)...)
	stdout, err := cmd.CombinedOutput()

	return common.CleanstdoutMultiline(stdout), err
}

func DeleteBranch(branch string) (string, error) {
	cmdLocal := exec.Command("git", "branch", "-D", branch)
	localStdout, err := cmdLocal.CombinedOutput()
	if err != nil {
		return common.CleanstdoutMultiline(localStdout), err
	}

	logging.Instance().Debugf("deleted local branch: %s", branch)

	cmdRemote := exec.Command("git", "push", "origin", "--delete", branch)
	remoteStdout, err := cmdRemote.CombinedOutput()
	if err != nil {
		return common.CleanstdoutMultiline(remoteStdout), err
	}

	logging.Instance().Debugf("deleted remote branch: %s", branch)

	return common.CleanstdoutMultiline(remoteStdout), nil
}

func CheckoutDefaultBranch() (string, error) {
	defaultBranch, err := OriginDefaultBranch()
	if err != nil {
		return defaultBranch, err
	}

	logging.Instance().Debugf("captured default branch: %s", defaultBranch)

	if stderr, err := Checkout(defaultBranch, false); err != nil {
		return stderr, err
	}

	logging.Instance().Debugf("checkout to default branch, %s, successful", defaultBranch)

	if stderr, err := PullChanges(); err != nil {
		return stderr, err
	}

	logging.Instance().Debug("pulled most recent changes for default branch")

	return "", nil
}

func StageChanges() (string, error) {
	cmd := exec.Command("git", "add", "-A")
	stderr, err := cmd.CombinedOutput()

	logging.Instance().Debug("successfully staged all changes for git repo")
	
	return common.CleanstdoutMultiline(stderr), err
}

func PushRemote(pushArgs string) (string, error) {
	pushCommand := fmt.Sprintf("git push %s", pushArgs)
	cmd := exec.Command("bash", "-c", pushCommand)
	stderr, err := cmd.CombinedOutput()

	return common.CleanstdoutMultiline(stderr), err
}

func GetPreviousNCommitMessage(count int) ([]string, error) {
	logging.Instance().Debug("capturing previous commits from git log")

	cmd := exec.Command("git", "log", "-" + fmt.Sprint(count), "--pretty=%B")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var commits []string
	for _, m := range strings.Split(string(stdout), "\n") {
		if m != "" {
			commits = append(commits, m)
		}
	}

	logging.Instance().Debugf("captured %d previous commits: %v", len(commits), commits)

	return commits, nil
}

func CreateTag(name, message string, force bool) (string, error) {
	var tagCmd *exec.Cmd
	if force {
		tagCmd = exec.Command("git", "tag", "-a", name, "--force", "-m", message)
	} else {
		tagCmd = exec.Command("git", "tag", "-a", name, "-m", message)
	}
	stdout, err := tagCmd.CombinedOutput()

	return common.CleanstdoutMultiline(stdout), err
}

func Rebase(branch string) (string, error) {
	cmd := exec.Command("git", "rebase", branch)
	stdout, err := cmd.CombinedOutput()
	
	return common.CleanstdoutMultiline(stdout), err
}

func LogFor(branch string) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git log --pretty=oneline --first-parent --format='`%%h` - %%s' | grep '%s'", branch))
	stdout, err := cmd.CombinedOutput()

	return common.CleanstdoutMultiline(stdout), err
}

func PushRemoteTagsOnly() (string, error) {
	cmd := exec.Command("git", "push", "--tags", "--force")
	stderr, err := cmd.CombinedOutput()
	
	return common.CleanstdoutMultiline(stderr), err
}

func LatestTagName() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return common.CleanStdoutSingleline(stdout), nil
}

func ProjectExistingVersionPrefix() (string, error) {
	tagName, err := LatestTagName()
	if err != nil {
		logging.Instance().Debugf("error ocurred when reading latest tagName from repo: %v\n%s", err, tagName)

		if strings.Contains(err.Error(), "128") {
			logging.Instance().Debug("existing tag prefix defaulting to global defaults since no existing tags found on remote")
			return config.AppConfig().TagPrefix(), nil
		}

		return "", fmt.Errorf("could not tag information from remote origin. %v", err)
	}

	var existingPrefix string
	if prefixSearch := regexp.MustCompile(constants.VersionPrefixRegexp).FindStringSubmatch(tagName); len(prefixSearch) > 0 {
		existingPrefix = strings.TrimSpace(prefixSearch[0])
	} else {
		existingPrefix = ""
	}

	logging.Instance().Debugf("captured existing prefix for repository: %s", existingPrefix)

	return existingPrefix, nil
}

func OriginCurrentVersion() (semver.Semver, error) {
	version := semver.Semver{0,0,0}

	defaultBranch, err := OriginDefaultBranch()
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