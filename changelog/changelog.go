package changelog

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"sykesdev.ca/gog/common"
	"sykesdev.ca/gog/logging"
	"sykesdev.ca/gog/models"
	"sykesdev.ca/gog/semver"
)

const changelogHeader = `
# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



`

func CreateChangeLogLines(entry *ChangelogEntry) ([]string, error) {
	workingDir, _ := common.WorkspacePaths()

	f, err := os.OpenFile(workingDir + "/CHANGELOG.md", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	logging.GetLogger().Debug("Creating changelog entry ...")

	var changelogLines []string
	scanner := bufio.NewScanner(bufio.NewReader(f))
	for scanner.Scan() {
		changelogLines = append(changelogLines, scanner.Text())
	}

	latestFeatIndex := 0
	versionLine := regexp.MustCompile(`^(#){2}(\ ){1}(\[)`)
	for i, line := range changelogLines {
		if matched := versionLine.MatchString(line); matched {
			latestFeatIndex = i
			break
		}
	}

	existingChanges := []string{}
	if latestFeatIndex != 0 {
		existingChanges = changelogLines[latestFeatIndex:]
	}

	changelogLines = append([]string{}, changelogHeader)
	changelogLines = append(changelogLines, entry.String())
	changelogLines = append(changelogLines, existingChanges...)

	return changelogLines, nil
}

func WriteChangelogToFile(lines []string) error {
	workingDir, _ := common.WorkspacePaths()

	changelogFile, err := os.Create(workingDir + "/CHANGELOG.md")
	if err != nil {
		return err
	}
	defer changelogFile.Close()
	
	_, err = changelogFile.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return err
	}

	return nil
}

type ChangelogEntry struct {
	Feature *models.Feature
	Version semver.Semver
	Added bool
}

func NewChangelogEntry(feature *models.Feature, version semver.Semver, added bool) (*ChangelogEntry) {
	return &ChangelogEntry{ Feature: feature, Version: version, Added: added }
}

func (e *ChangelogEntry) Lines() []string {
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

	changes, err := e.Feature.ListChanges()
	if err != nil {
		logging.GetLogger().Fatal(fmt.Sprintf("failed to get feature changes from git. try pushing a change first. %v", err))
	}

	lines = append(lines, changes...)
	lines = append(lines, "\n\n")

	return lines
}

func (e *ChangelogEntry) String() string {
	return strings.Join(e.Lines(), "\n")
}

