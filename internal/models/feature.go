package models

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/common/constants"
	"sykesdev.ca/gog/internal/git"
	"sykesdev.ca/gog/internal/semver"
)

type Feature struct {
	Jira string `json:"jira"`
	Comment string `json:"comment"`
	CustomVersionPrefix string `json:"custom_prefix"`
	TestCount int `json:"test_count"`
}

func NewFeature(jira, comment, versionPrefix string) (*Feature, error) {
	feat := &Feature{Jira: jira, Comment: comment, TestCount: 0}

	if versionPrefix != "" {
		if matched, _ := regexp.MatchString(constants.VersionPrefixRegexp, versionPrefix); !matched {
			return nil, errors.New("invalid version prefix specified for feature")
		}
		feat.CustomVersionPrefix = versionPrefix
	}

	return feat, nil
}

func NewFeatureFromFile() (*Feature, error) {
	GOGDir := common.GOGPath()
	
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

func (f *Feature) CreateBranch() (string, error) {
	return git.Checkout(f.Jira, true)
}

func (f *Feature) DeleteBranch() (string, error) {
	if cbStdout, err := git.GetCurrentBranch(); cbStdout == f.Jira {
		git.CheckoutDefaultBranch()
	} else if err != nil {
		return cbStdout, err
	}

	return git.DeleteBranch(f.Jira)
}

func (f *Feature) LocalExists() bool {
	return git.LocalBranchExists(f.Jira)
}

func (f *Feature) RemoteExists() bool {
	return git.RemoteBranchExists(f.Jira)
}

func (f *Feature) PushChanges(commitMessage string) (string, error) {
	if stderr, err := git.StageChanges(); err != nil {
		return stderr, err
	}

	var pushArgs string
	if git.HasUncommittedChanges() {
		if stderr, err := f.CommitChanges(commitMessage); err != nil {
			return stderr, err
		}

		if !f.RemoteExists() {
			pushArgs = fmt.Sprintf("--set-upstream origin %s", f.Jira)
		} else {
			// only pull changes if a remote exists
			if stderr, err := git.PullChanges(); err != nil {
				return stderr, err
			}
		}
	}

	return git.PushRemote(pushArgs)
}

func (f *Feature) CreateReleaseTags(version semver.Semver) (string, error) {
	tagMessage := fmt.Sprintf("(%s): %s %s", version, f.Jira, f.Comment)
	tagStdout, err := git.CreateTag(version.String(), tagMessage, false)
	if err != nil {
		return tagStdout, err
	}

	tagMessage = fmt.Sprintf("(%s): %s %s", version.Major(), f.Jira, f.Comment)
	tagStdout, err = git.CreateTag(version.Major(), tagMessage, true)
	
	return tagStdout, err
}

func (f *Feature) Rebase() (string, error) {
	return git.Rebase(f.Jira)
}

func (f *Feature) ListFeatureChanges() ([]string, error) {
	var changes []string
	changeBlob, err := git.LogFor(f.Jira)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(changeBlob)))
	for scanner.Scan() {
		changes = append(changes, fmt.Sprintf("- %s", scanner.Text()))
	}

	return changes, nil
}

func (f *Feature) CommitChanges(commitMessage string) (string, error) {
	formattedMessage := fmt.Sprintf("%s %s", f.Jira, commitMessage)

	return git.Commit(formattedMessage)
}

func (f *Feature) Save() error {
	GOGDir := common.GOGPath()

	if !common.PathExists(GOGDir) {
		if err := os.MkdirAll(GOGDir, 0700); err != nil {
			return err
		}
	}

	featureFile, err := os.Create(GOGDir + "/feature.json")
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

func (f *Feature) Clean() error {
	GOGDir := common.GOGPath()

	if _, err := git.CheckoutDefaultBranch(); err != nil {
		return err
	}

	if _, err := f.DeleteBranch(); err != nil {
		return err
	}

	if err := os.RemoveAll(GOGDir); err != nil {
		return err
	}

	return nil
}