package lib

type Feature struct {
	Jira string `json:"jira"`
	Comment string `json:"comment"`
}

func NewFeature(jira, comment string) (*Feature, error) {
	feat := &Feature{Jira: jira, Comment: comment}

	return feat, nil
}