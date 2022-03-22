package changelog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"sykesdev.ca/gog/lib"
	"sykesdev.ca/gog/lib/semver"
)

func CreateChangeLogLines(entry *ChangelogEntry) ([]string, error) {
	workingDir, _ := lib.WorkspacePaths()

	f, err := os.OpenFile(workingDir + "/CHANGELOG.md", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	lib.GetLogger().Debug("Creating changelog entry ...")

	var changelogLines []string
	scanner := bufio.NewScanner(bufio.NewReader(f))
	for scanner.Scan() {
		changelogLines = append(changelogLines, scanner.Text())
	}

	if len(changelogLines) > 0 {
		lib.GetLogger().Debug("CHANGELOG.md already exists ... will append new changelog entry")

		var latestFeatIndex int
		for i, line := range changelogLines {
			if strings.Contains(line, "## [") {
				latestFeatIndex = i
				break
			}
		}

		existingChanges := changelogLines[latestFeatIndex:]
		changelogLines = append(changelogLines[:latestFeatIndex-1], entry.String())
		changelogLines = append(changelogLines, existingChanges...)

		return changelogLines, nil
	} else {
		changelogLines = append(changelogLines,
`
# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



`)
		changelogLines = append(changelogLines, entry.String())
	}

	return changelogLines, nil
}

func WriteChangelogToFile(lines []string) error {
	workingDir, _ := lib.WorkspacePaths()

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
	Feature *lib.Feature
	Version semver.Semver
	Added bool
}

func NewChangelogEntry(feature *lib.Feature, version semver.Semver, added bool) (*ChangelogEntry) {
	return &ChangelogEntry{ Feature: feature, Version: version, Added: added }
}

func (e *ChangelogEntry) String() string {
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
	lines = append(lines, fmt.Sprintf("\n> %s =-= %s", e.Feature.Jira, e.Feature.Comment))

	if e.Added {
		lines = append(lines, "\n### Added\n")
	} else {
		lines = append(lines, "\n### Changed\n")
	}

	changes, err := e.Feature.ListChanges()
	if err != nil {
		lib.GetLogger().Fatal(fmt.Sprintf("Failed to get feature changes from git. %v", err))
	}

	lines = append(lines, changes...)
	lines = append(lines, "\n\n")
	
	return strings.Join(lines, "\n")
}

