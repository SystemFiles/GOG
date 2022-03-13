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
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

func GitPullChanges() error {
	cmd := exec.Command("git", "pull")
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func GitCreateBranch(name string, checkout bool) error {
	if checkout {
		cmd := exec.Command("git", "checkout", "-b", name)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command("git", "branch", name)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}

	return nil
}

func GitOriginDefaultBranch() (string, error) {
	defaultBranchCmd := exec.Command("bash", "-c", "git remote show origin | sed -n '/HEAD branch/s/.*: //p'")
	defaultBranch, err := defaultBranchCmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(defaultBranch)), nil
}

func GitCheckoutDefaultBranch() error {
	defaultBranch, err := GitOriginDefaultBranch()
	if err != nil {
		return err
	}

	checkoutCmd := exec.Command("git", "checkout", defaultBranch)
	_, err = checkoutCmd.Output()
	if err != nil {
		return err
	}

	if err = GitPullChanges(); err != nil {
		return err
	}

	return nil
}

func GitStageChanges() error {
	cmd := exec.Command("git", "add", "-A")
	_, err := cmd.Output()
	
	return err
}

func GitCommitChanges(feature *Feature, commitMessage string) error {
	formattedMessage := fmt.Sprintf("%s %s", feature.Jira, commitMessage)
	cmd := exec.Command("git", "commit", "-m", formattedMessage)
	_, err := cmd.Output()
	
	return err
}

func GitRemoteExists(feature *Feature) bool {
	remoteExistsCommand := fmt.Sprintf("git ls-remote --heads --exit-code | egrep %s", feature.Jira)
	cmd := exec.Command("bash", "-c", remoteExistsCommand)
	_, err := cmd.Output()

	return err == nil
}

func GitPushRemote(pushArgs string) error {
	pushCommand := fmt.Sprintf("git push %s", pushArgs)
	cmd := exec.Command("bash", "-c", pushCommand)
	_, err := cmd.Output()

	return err
}