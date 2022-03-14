package lib

import (
	"fmt"
	"os/exec"
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