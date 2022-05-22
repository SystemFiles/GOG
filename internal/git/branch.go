package git

import (
	"fmt"
	"os/exec"

	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/logging"
)

type Branch struct {
	Name string `json:"name"`
	RemoteExists bool `json:"remote_exists"`
	LocalExists bool `json:"local_exists"`
}

func NewBranch(name string) *Branch {
	b := &Branch{
		Name: name,
	}
	b.RemoteExists = remoteBranchExists(b)
	b.LocalExists = localBranchExists(b)

	return b
}

func (b *Branch) UpdateBranch(branch string) {
	b.Name = branch

	b.RemoteExists = remoteBranchExists(b)
	b.LocalExists = localBranchExists(b)
}

func (b *Branch) UncommittedChanges() bool {
	cmd := exec.Command("bash", "-c", "git status --porcelain | egrep '^[A,M,D,R]'")
	_, err := cmd.Output()

	logging.Instance().Debugf("uncommitted changes: %t", err == nil)

	return err == nil
}

func (b *Branch) RelatedLogs() (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git log --pretty=oneline --first-parent --format='`%%h` - %%s' | grep '%s'", b.Name))
	stdout, err := cmd.CombinedOutput()

	return common.CleanstdoutMultiline(stdout), err
}

func (b *Branch) Delete() (error) {
	cmdLocal := exec.Command("git", "branch", "-D", b.Name)
	localStdout, err := cmdLocal.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(localStdout))
	}

	logging.Instance().Debugf("deleted local branch: %s", b.Name)

	cmdRemote := exec.Command("git", "push", "origin", "--delete", b.Name)
	remoteStdout, err := cmdRemote.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(remoteStdout))
	}

	logging.Instance().Debugf("deleted remote branch: %s", b.Name)

	return nil
}

func (b *Branch) String() string {
	return fmt.Sprintf("Branch: { Name: %s, Remote Exists: %t }", b.Name, b.RemoteExists)
}