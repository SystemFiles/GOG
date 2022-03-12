package lib

type Feature struct {
	Jira string `yaml:"jira"`
	Comment string `yaml:"comment"`
}

func NewFeature(jira, comment string) (*Feature, error) {
	feat := &Feature{Jira: jira, Comment: comment}

	return feat, nil
}