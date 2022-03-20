package lib

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func CreateChangeLogLines(entry *ChangelogEntry) ([]string, error) {
	workingDir, _ := GetWorkspacePaths()

	f, err := os.Create(workingDir + "/CHANGELOG.md")
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	var lines []byte
	_, err = f.Read(lines)
	if err != nil {
		return []string{}, err
	}

	fmt.Println(*entry)

	changlogLines := strings.Split(string(lines), "\n")

	fmt.Println(changlogLines, len(changlogLines))

	if (len(changlogLines) - 1) > 0 {
		GetLogger().Info("CHANGELOG already exists")

		for _, line := range changlogLines {
			if strings.Contains(line, "## [") {
				GetLogger().Debug("Found start of existing CHANGELOG entry")
			}
		}
	} else {
		changlogLines = append(changlogLines,
`
# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



`)

		GetLogger().Info("Creating changelog")
		changlogLines = append(changlogLines, entry.String())
	}

	return changlogLines, nil
}

func WriteChangelogToFile(lines []string) error {
	workingDir, _ := GetWorkspacePaths()

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
	Feature *Feature
	Version string
	Added bool
}

func NewChangelogEntry(feature *Feature, version string, added bool) (*ChangelogEntry) {
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

	changes, err := GitViewFeatureChanges(e.Feature)
	if err != nil {
		panic("Failed to get feature changes from git")
	}

	lines = append(lines, changes...)
	
	return strings.Join(lines, "\n")
}

