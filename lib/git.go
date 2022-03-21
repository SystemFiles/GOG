package lib

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func GitHasUnstagedCommits() bool {
	cmd := exec.Command("git", "diff-index", "--quiet", "HEAD")
	_, err := cmd.Output()

	return err != nil
}

func GitHasUncommittedChanges() bool {
	cmd := exec.Command("bash", "-c", "git status --porcelain | egrep '^[A,M,D]'")
	_, err := cmd.Output()

	return err == nil
}

func GitIsValidRepo() bool {
	cmd := exec.Command("git", "status")
	_, err := cmd.Output()
	
	return err == nil
}

func GitBranchExists(branch string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git branch | egrep %s", branch))
	_, err := cmd.Output()

	return err == nil
}

func GitGetCurrentBranch() (string, error) {
	cmd := exec.Command("bash", "-c", "git branch | grep '*' | cut -d' ' -f2")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

func GitPullChanges() (string, error) {
	cmd := exec.Command("git", "pull")
	stdout, err := cmd.CombinedOutput()
	
	return string(stdout), err
}

func GitOriginDefaultBranch() (string, error) {
	defaultBranchCmd := exec.Command("bash", "-c", "git remote show origin | sed -n '/HEAD branch/s/.*: //p'")
	defaultBranch, err := defaultBranchCmd.CombinedOutput()

	return strings.TrimSpace(string(defaultBranch)), err
}

func GitCheckoutDefaultBranch() (string, error) {
	defaultBranch, err := GitOriginDefaultBranch()
	if err != nil {
		return defaultBranch, err
	}

	checkoutCmd := exec.Command("git", "checkout", defaultBranch)
	checkoutStdout, err := checkoutCmd.CombinedOutput()
	if err != nil {
		return string(checkoutStdout), err
	}

	if stderr, err := GitPullChanges(); err != nil {
		return stderr, err
	}

	return "", nil
}

func GitStageChanges() (string, error) {
	cmd := exec.Command("git", "add", "-A")
	stderr, err := cmd.CombinedOutput()
	
	return string(stderr), err
}

func GitCommitChanges(feature *Feature, commitMessage string) (string, error) {
	formattedMessage := fmt.Sprintf("%s %s", feature.Jira, commitMessage)
	cmd := exec.Command("git", "commit", "-m", formattedMessage)
	stderr, err := cmd.CombinedOutput()
	
	return string(stderr), err
}

func GitPushRemote(pushArgs string) (string, error) {
	pushCommand := fmt.Sprintf("git push %s", pushArgs)
	cmd := exec.Command("bash", "-c", pushCommand)
	stderr, err := cmd.CombinedOutput()

	return string(stderr), err
}

func GitPushRemoteTagsOnly() (string, error) {
	cmd := exec.Command("git", "push", "--tags", "--force")
	stderr, err := cmd.CombinedOutput()
	
	return string(stderr), err
}

func GitFeatureChanges(feature *Feature) ([]string, error) {
	var changes []string
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git log --pretty=oneline --first-parent --format='`%%h` - %%s' | grep '%s'", feature.Jira))
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	for scanner.Scan() {
		changes = append(changes, fmt.Sprintf("- %s", scanner.Text()))
	}

	return changes, nil
}

func GitOriginCurrentVersion() ([3]int, error) {
	version := [3]int{0,0,0}

	defaultBranch, err := GitOriginDefaultBranch()
	if err != nil {
		return version, err
	}

	tagCmd := exec.Command("bash", "-c", fmt.Sprintf("git tag --merged %s --sort=taggerdate | tail -r", defaultBranch))
	tagOut, err := tagCmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "128") {
			return version, nil
		}
		
		return version, err
	}

	semverRegex, err := regexp.Compile(`^([0-9])+\.([0-9])+\.([0-9])$`)
	if err != nil {
		return version, err
	}

	var latestTag string
	tagScanner := bufio.NewScanner(strings.NewReader(string(tagOut)))
	for tagScanner.Scan() {
		tag := tagScanner.Text()
		if matched := semverRegex.MatchString(tag); matched {
			latestTag = tag
			break
		}
	}

	if latestTag == "" {
		return version, nil
	}

	verElements := strings.Split(string(latestTag), ".")
	major, err := strconv.Atoi(verElements[0])
	if err != nil { return version, err }
	minor, err := strconv.Atoi(verElements[1])
	if err != nil { return version, err }
	patch, err := strconv.Atoi(verElements[2])
	if err != nil { return version, err }

	version = [3]int{major, minor, patch}

	return version, nil
}

func GitRebase(feature *Feature) (string, error) {
	cmd := exec.Command("git", "rebase", feature.Jira)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), err
	}

	return string(stdout), nil
}

func GitCreateReleaseTags(version string, feature *Feature) (string, error) {
	tagCmd := exec.Command("git", "tag", "-a", version, "-m", fmt.Sprintf("(%s): %s %s", version, feature.Jira, feature.Comment))
	stdout, err := tagCmd.CombinedOutput()
	if err != nil {
		return string(stdout), err
	}

	majorVersion := strings.Split(version, ".")[0] + ".x"
	majorTagCmd := exec.Command("git", "tag", "-a", majorVersion, "--force", "-m", fmt.Sprintf("(%s): %s %s", majorVersion, feature.Jira, feature.Comment))
	stdout, err = majorTagCmd.CombinedOutput()
	
	return string(stdout), err
}