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

func (f *Feature) CreateReleaseTags(r *git.Repository, version semver.Semver) error {
	tagMessage := fmt.Sprintf("(%s): %s %s", version, f.Jira, f.Comment)

	logging.Instance().Debugf("creating release tag with message: %s", tagMessage)

	err := r.CreateTag(version.String(), tagMessage, false)
	if err != nil {
		return err
	}

	logging.Instance().Debugf("created release tag (%s) for feature: %s", version, f)

	tagMessage = fmt.Sprintf("(%s): %s %s", version.Major(), f.Jira, f.Comment)
	err = r.CreateTag(version.Major(), tagMessage, true)

	logging.Instance().Debugf("created release tag (%s) for feature: %s", version.Major(), f)
	
	return err
}

func (f *Feature) Changes(r *git.Repository) ([]string, error) {
	var changes []string
	changeBlob, err := r.FeatureBranch.RelatedLogs()
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

func (f *Feature) String() string {
	return fmt.Sprintf("%s %s", f.Jira, f.Comment)
}