package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Feature struct {
	Jira string `json:"jira"`
	Comment string `json:"comment"`
	TestCount int64 `json:"test_count"`
}

func NewFeature(jira, comment string) (*Feature, error) {
	feat := &Feature{Jira: jira, Comment: comment, TestCount: 0}

	return feat, nil
}

func NewFeatureFromFile() (*Feature, error) {
	_, GOGDir := GetWorkspacePaths()
	
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

func (f *Feature) CreateBranch(checkout bool) (string, error) {
	var stdout []byte
	var err error
	if checkout {
		cmd := exec.Command("git", "checkout", "-b", f.Jira)
		stdout, err = cmd.CombinedOutput()
	} else {
		cmd := exec.Command("git", "branch", f.Jira)
		stdout, err = cmd.CombinedOutput()
	}

	return string(stdout), err
}

func (f *Feature) RemoteExists() bool {
	remoteExistsCommand := fmt.Sprintf("git ls-remote --heads --exit-code | egrep %s", f.Jira)
	cmd := exec.Command("bash", "-c", remoteExistsCommand)
	_, err := cmd.Output()

	return err == nil
}



func (f *Feature) Save() error {
	workingDir, GOGDir := GetWorkspacePaths()

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