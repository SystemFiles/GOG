package changelog

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/git"
	"sykesdev.ca/gog/internal/logging"
	"sykesdev.ca/gog/internal/models"
	"sykesdev.ca/gog/internal/semver"
)

const fileHeader = `
# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



`

func CreateChangeLogLines(entry *ChangelogEntry) ([]string, error) {
	projectRoot, err := common.GitProjectRoot()
	if err != nil {
		return nil, err
	}

	logging.Instance().Debugf("projectRoot: %s", projectRoot)
	logging.Instance().Debugf("opening changelog file if exists from %s", projectRoot + "/CHANGELOG.md")

	f, err := os.OpenFile(projectRoot + "/CHANGELOG.md", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	logging.Instance().Debug("creating the changelog entry ...")

	var changelogLines []string
	scanner := bufio.NewScanner(bufio.NewReader(f))
	for scanner.Scan() {
		logging.Instance().Debugf("adding line from CHANGELOG.md: %s", scanner.Text())
		changelogLines = append(changelogLines, scanner.Text())
	}

	latestFeatIndex := 0
	versionLine := regexp.MustCompile(`^(#){2}(\ ){1}(\[)`)
	for i, line := range changelogLines {
		logging.Instance().Debugf("processing changelog line: %s", line)
		if matched := versionLine.MatchString(line); matched {
			logging.Instance().Debugf("matched version line: %s", line)
			latestFeatIndex = i
			break
		}
	}

	existingChanges := []string{}
	if latestFeatIndex != 0 {
		existingChanges = changelogLines[latestFeatIndex:]
	}

	logging.Instance().Debugf("captured existing changes with line-count: %d", len(existingChanges))
	logging.Instance().Debug("inserting changelog lines for new release ...")

	changelogLines = append([]string{}, fileHeader)
	changelogLines = append(changelogLines, entry.String())
	changelogLines = append(changelogLines, existingChanges...)

	logging.Instance().Debug("finished inserting changelog lines for new release")

	return changelogLines, nil
}

func WriteChangelogToFile(lines []string) error {
	projectRoot, err := common.GitProjectRoot()
	if err != nil {
		return err
	}

	logging.Instance().Debugf("projectRoot: %s", projectRoot)
	logging.Instance().Debugf("opening CHANGELOG for writing from: %s", projectRoot + "/CHANGELOG.md")

	changelogFile, err := os.Create(projectRoot + "/CHANGELOG.md")
	if err != nil {
		return err
	}
	defer changelogFile.Close()
	
	logging.Instance().Debugf("writing %d lines to CHANGELOG.md", len(lines))

	_, err = changelogFile.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return err
	}

	logging.Instance().Debug("completed write to CHANGELOG.md")

	return nil
}

type ChangelogEntry struct {
	Repository *git.Repository
	Feature *models.Feature
	Version semver.Semver
	Added bool
}

func NewChangelogEntry(feature *models.Feature, repo *git.Repository, version semver.Semver, added bool) (*ChangelogEntry) {
	return &ChangelogEntry{ Feature: feature, Repository: repo, Version: version, Added: added }
}

func (e *ChangelogEntry) Lines() []string {
	logging.Instance().Debug("generating lines for configured changelog entry")

	currentTime := time.Now().UTC()
	formattedTimeString := fmt.Sprintf("%d-%d-%d %d:%d:%d",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())
	var lines []string
	lines = append(lines, fmt.Sprintf("## [ %s ] - %s", e.Version, formattedTimeString))
	lines = append(lines, fmt.Sprintf("\n> %s %s", e.Feature.Jira, e.Feature.Comment))

	if e.Added {
		lines = append(lines, "\n### Added\n")
	} else {
		lines = append(lines, "\n### Changed\n")
	}

	changes, err := e.Feature.Changes(e.Repository)
	if err != nil {
		logging.Instance().Fatalf("failed to get feature changes from git. try pushing a change first. %v", err)
	}

	logging.Instance().Debugf("captured the following changes for this feature release: %v", changes)

	lines = append(lines, changes...)
	lines = append(lines, "\n\n")

	return lines
}

func (e *ChangelogEntry) String() string {
	return strings.Join(e.Lines(), "\n")
}

