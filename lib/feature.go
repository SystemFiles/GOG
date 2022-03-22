package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"sykesdev.ca/gog/lib/semver"
)

type Feature struct {
	Jira string `json:"jira"`
	Comment string `json:"comment"`
	TestCount int `json:"test_count"`
}

func NewFeature(jira, comment string) (*Feature, error) {
	feat := &Feature{Jira: jira, Comment: comment, TestCount: 0}

	return feat, nil
}

func NewFeatureFromFile() (*Feature, error) {
	_, GOGDir := WorkspacePaths()
	
	featureBytes, err := os.ReadFile(GOGDir + "/feature.json")
	if err != nil {
		return nil, err
	}

	var feature *Feature
	err = json.Unmarshal(featureBytes, &feature)
	if err != nil {
		return nil, err
	}

	return feature, nil
}

func (f *Feature) UpdateTestCount() error {
	f.TestCount += 1
	
	if err := f.Save(); err != nil {
		return err
	}

	return nil
}

func (f *Feature) BranchExists() bool {
	return GitBranchExists(f.Jira)
}

func (f *Feature) CreateBranch() (string, error) {
	cmd := exec.Command("git", "checkout", "-b", f.Jira)
	stdout, err := cmd.CombinedOutput()

	return string(stdout), err
}

func (f *Feature) DeleteBranch() (string, error) {
	if cbStdout, err := GitGetCurrentBranch(); cbStdout == f.Jira {
		GitCheckoutDefaultBranch()
	} else if err != nil {
		return cbStdout, err
	}

	cmdLocal := exec.Command("git", "branch", "-D", f.Jira)
	localStdout, err := cmdLocal.CombinedOutput()
	if err != nil {
		return string(localStdout), err
	}

	cmdRemote := exec.Command("git", "push", "origin", "--delete", f.Jira)
	remoteStdout, err := cmdRemote.CombinedOutput()
	
	return string(remoteStdout), err
}

func (f *Feature) RemoteExists() bool {
	remoteExistsCommand := fmt.Sprintf("git ls-remote --heads --exit-code | egrep %s", f.Jira)
	cmd := exec.Command("bash", "-c", remoteExistsCommand)
	_, err := cmd.Output()

	return err == nil
}

func (f *Feature) PushChanges(commitMessage string) (string, error) {
	if stderr, err := GitStageChanges(); err != nil {
		return stderr, err
	}

	if stderr, err := GitCommitChanges(f, commitMessage); err != nil {
		return stderr, err
	}

	var pushArgs string
	if !f.RemoteExists() {
		pushArgs = fmt.Sprintf("--set-upstream origin %s", f.Jira)
	} else {
		// only pull changes if a remote exists
		if stderr, err := GitPullChanges(); err != nil {
			return stderr, err
		}
	}

	stderr, err := GitPushRemote(pushArgs)

	return stderr, err
}

func (f *Feature) CreateReleaseTags(version semver.Semver) (string, error) {
	tagCmd := exec.Command("git", "tag", "-a", version.String(), "-m", fmt.Sprintf("(%s): %s %s", version, f.Jira, f.Comment))
	stdout, err := tagCmd.CombinedOutput()
	if err != nil {
		return string(stdout), err
	}

	majorTagCmd := exec.Command("git", "tag", "-a", version.Major(), "--force", "-m", fmt.Sprintf("(%s): %s %s", version.Major(), f.Jira, f.Comment))
	stdout, err = majorTagCmd.CombinedOutput()
	
	return string(stdout), err
}

func (f *Feature) Rebase() (string, error) {
	cmd := exec.Command("git", "rebase", f.Jira)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), err
	}

	return string(stdout), nil
}

func (f *Feature) ListChanges() ([]string, error) {
	var changes []string
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git log --pretty=oneline --first-parent --format='`%%h` - %%s' | grep '%s'", f.Jira))
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

func (f *Feature) Save() error {
	workingDir, GOGDir := WorkspacePaths()

	if !PathExists(GOGDir) {
		if err := os.MkdirAll(GOGDir, 0700); err != nil {
			return err
		}
	}

	featureFile, err := os.Create(workingDir + "/.gog/feature.json")
	if err != nil {
		return err
	}
	defer featureFile.Close()

	featureBytes, err := json.Marshal(f)
	if err != nil {
		return err
	}
	
	_, err = featureFile.Write(featureBytes)
	if err != nil {
		return err
	}

	return nil
}