package lib

import (
	"encoding/json"
	"os"
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

func (f *Feature) UpdateTestCount() error {
	f.TestCount += 1
	
	if err := f.Save(); err != nil {
		return err
	}

	return nil
}

func (f *Feature) Save() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
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