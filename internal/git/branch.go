package git

import (
	"fmt"
	"os/exec"

	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/logging"
)

// TODO should not transform its own state - remote update and find alternative workflow
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

func (b *Branch) String() string {
	return b.Name
}