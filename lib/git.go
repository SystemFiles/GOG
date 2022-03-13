package lib

import "os/exec"

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