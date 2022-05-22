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
	"sykesdev.ca/gog/internal/logging"
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

	logging.Instance().Debugf("created new feature: %s", feat)

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

	logging.Instance().Debugf("created feature instance from file: %s", feature)

	return feature, nil
}

func (f *Feature) UpdateTestCount() error {
	f.TestCount += 1
	
	if err := f.Save(); err != nil {
		return err
	}

	logging.Instance().Debugf("updated test build count from %d -> %d", f.TestCount - 1, f.TestCount)

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

	logging.Instance().Debugf("deleting feature branch: %s", f.Jira)

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

		logging.Instance().Debug("feature has uncommitted changes. committing them now")

		if stderr, err := f.CommitChanges(commitMessage); err != nil {
			return stderr, err
		}

		logging.Instance().Debug("changes committed. pushing changes to remote")

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

	logging.Instance().Debugf("created release tag (%s) for feature: %s", version, f)

	tagMessage = fmt.Sprintf("(%s): %s %s", version.Major(), f.Jira, f.Comment)
	tagStdout, err = git.CreateTag(version.Major(), tagMessage, true)

	logging.Instance().Debugf("created release tag (%s) for feature: %s", version.Major(), f)
	
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

	logging.Instance().Debugf("got change blob from git logs containing: %s", changeBlob)
	logging.Instance().Debug("creating formatted change log entries for feature")

	scanner := bufio.NewScanner(strings.NewReader(string(changeBlob)))
	for scanner.Scan() {
		changes = append(changes, fmt.Sprintf("- %s", scanner.Text()))
	}

	logging.Instance().Debugf("feature changes captured: %v", changes)

	return changes, nil
}

func (f *Feature) CommitChanges(commitMessage string) (string, error) {
	formattedMessage := fmt.Sprintf("%s %s", f.Jira, commitMessage)

	logging.Instance().Debugf("committing changes with message: %s", formattedMessage)

	return git.Commit(formattedMessage)
}

func (f *Feature) Save() error {
	GOGDir := common.GOGPath()

	if !common.PathExists(GOGDir) {
		if err := os.MkdirAll(GOGDir, 0700); err != nil {
			return err
		}
	}

	logging.Instance().Debugf("saving feature changes to file at: %s", GOGDir + "/feature.json")

	featureFile, err := os.Create(GOGDir + "/feature.json")
	if err != nil {
		return err
	}
	defer featureFile.Close()

	featureBytes, err := json.Marshal(f)
	if err != nil {
		return err
	}

	logging.Instance().Debugf("serialized feature into %d bytes", len(featureBytes))
	
	_, err = featureFile.Write(featureBytes)
	if err != nil {
		return err
	}

	logging.Instance().Debug("successfully wrote feature changes to file")

	return nil
}

func (f *Feature) Clean() error {
	GOGDir := common.GOGPath()

	logging.Instance().Debug("cleaning feature files")

	if _, err := git.CheckoutDefaultBranch(); err != nil {
		return err
	}

	if _, err := f.DeleteBranch(); err != nil {
		return err
	}

	logging.Instance().Debugf("deleted feature branch: %s", f.Jira)

	if err := os.RemoveAll(GOGDir); err != nil {
		return err
	}

	logging.Instance().Debug("removed all GOG feature files from project")

	return nil
}

func (f *Feature) String() string {
	return fmt.Sprintf("%s %s", f.Jira, f.Comment)
}