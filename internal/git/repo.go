package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/logging"
	"sykesdev.ca/gog/internal/semver"
)

type Repository struct {
	Name string
	VersionPrefix string

	DefaultBranch *Branch
	FeatureBranch *Branch
	CurrentBranch *Branch
	LastTag semver.Semver

	mutex sync.Mutex
}

func NewRepository() (*Repository, error) {
	var wg sync.WaitGroup
	r := &Repository{mutex: sync.Mutex{}, FeatureBranch: &Branch{}}

	rootChan := make(chan []string)
	prefixChan, dBranchChan, cBranchChan := make(chan string), make(chan string), make(chan string)
	lastTagChan := make(chan semver.Semver)
	doneChan := make(chan struct{})
	errChan := make(chan error)

	if !repositoryIsValid() {
		return nil, errors.New("directory does not contain a valid git repository")
	}

	wg.Add(1)
	go func () {
		root, err := common.GitProjectRoot()
		if err != nil {
			errChan <- err
		}
		rootChan <- strings.Split(root, "/")

		logging.Instance().Debugf("completed search for project root with result: %s", root)

		wg.Done()
	}()
	
	wg.Add(1)
	go func ()  {
		existingPrefix, err := projectExistingVersionPrefix()
		if err != nil {
			errChan <- err
		}

		prefixChan <- existingPrefix

		logging.Instance().Debugf("completed prefix with result: %s", existingPrefix)

		wg.Done()
	}()

	wg.Add(1)
	go func ()  {
		defaultBranch, err := originDefaultBranch()
		if err != nil {
			errChan <- err
		}

		dBranchChan <- defaultBranch

		logging.Instance().Debugf("completed default branch with result: %s", defaultBranch)

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		currentBranch, err := getCurrentBranch()
		if err != nil {
			errChan <- err
		}

		cBranchChan <- currentBranch

		logging.Instance().Debugf("completed current branch with result: %s", currentBranch)

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		latestTag, err := originLatestFullVersion()
		if err != nil {
			errChan <- err
		}

		lastTagChan <- latestTag

		logging.Instance().Debugf("completed latest tag search with result: %s", latestTag)

		wg.Done()
	}()

	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	for {
		select {
			case rootParts := <- rootChan:
				r.Name = strings.TrimSpace(rootParts[len(rootParts) - 1])
			case existingPrefix := <- prefixChan:
				r.VersionPrefix = existingPrefix
			case defaultBranch := <- dBranchChan:
				r.DefaultBranch = NewBranch(defaultBranch)
			case currentBranch := <- cBranchChan:
				r.CurrentBranch = NewBranch(currentBranch)
			case tagName := <- lastTagChan:
				r.LastTag = tagName
			case e := <- errChan:
				logging.Instance().Debugf("error ocurred when capturing repository metadata. %v", e)
				return nil, e
			case <- doneChan:
				logging.Instance().Debugf("initialized repository with values %v", r)

				if r.DefaultBranch == nil || r.CurrentBranch == nil || r.Name == "" {
					return nil, errors.New("failed to initialize GOG repository")
				}

				return r, nil
		}
	}
}

func (r *Repository) ContainsBranch(branch string) bool {
	b := NewBranch(branch)
	return b.RemoteExists || b.LocalExists
}

func (r *Repository) CheckoutBranch(branch *Branch, create, isFeature bool) error {
	checkoutArgs := make([]string, 0)
	if create {
		checkoutArgs = append(checkoutArgs, "-b")
	}
	checkoutArgs = append(checkoutArgs, branch.Name)

	logging.Instance().Debugf("checking out branch, %s, with create: %t", branch, create)

	cmd := exec.Command("git", append([]string{"checkout"}, checkoutArgs...)...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stdout))
	}

	r.CurrentBranch.UpdateBranch(branch.Name)
	if create && isFeature {
		r.FeatureBranch = NewBranch(branch.Name)
	}

	return nil
}

func (r *Repository) DeleteFeatureBranch() error {
	return deleteBranch(r.FeatureBranch)
}

func (r *Repository) StageChanges() error {
	cmd := exec.Command("git", "add", "-A")
	stderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stderr))
	}

	logging.Instance().Debug("successfully staged all changes for git repo")

	return nil
}

func (r *Repository) CommitChanges(message string) error {
	if r.CurrentBranch.UncommittedChanges() {
		logging.Instance().Debugf("uncommitted changes found... committing them with message: %s", message)
		cmd := exec.Command("git", "commit", "-m", message)
		stderr, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stderr))
		}
	}

	return nil
}

func (r *Repository) PullChanges() error {
	cmd := exec.Command("git", "fetch", "--tags", "--force")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stdout))
	}

	logging.Instance().Debugf("fetched tags from remote with output: %s", string(stdout))

	cmd = exec.Command("git", "pull", "--all")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stdout))
	}
	
	logging.Instance().Debugf("pulled changes from remote with output: %s", string(stdout))

	return nil
}

func (r *Repository) Push() error {
	var pushArgs string
	if !r.CurrentBranch.RemoteExists {
		pushArgs = "--set-upstream origin " + r.CurrentBranch.Name
	}

	logging.Instance().Debugf("pushing changes for %s with the following arguments: %s", r.CurrentBranch.Name, pushArgs)

	pushCommand := fmt.Sprintf("git push %s", pushArgs)
	cmd := exec.Command("bash", "-c", pushCommand)
	stderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stderr))
	}

	return nil
}

func (r *Repository) PushTags() error {
	cmd := exec.Command("git", "push", "--tags", "--force")
	stderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stderr))
	}
	
	return nil
}

func (r *Repository) Rebase() error {
	cmd := exec.Command("git", "rebase", r.DefaultBranch.Name)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stdout))
	}

	return nil
}

func (r *Repository) SquashMerge() error {
	cmd := exec.Command("git", "merge", "--squash", r.FeatureBranch.Name)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stdout))
	}

	return nil
}

func (r *Repository) LogN(N int) ([]string, error) {
	logging.Instance().Debugf("capturing previous %d commits from git log", N)

	cmd := exec.Command("git", "log", "-" + fmt.Sprint(N), "--pretty=%B")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var commits []string
	for _, m := range strings.Split(string(stdout), "\n") {
		if m != "" {
			commits = append(commits, m)
		}
	}

	logging.Instance().Debugf("captured %d previous commits: %v", len(commits), commits)

	return commits, nil
}

func (r *Repository) CreateTag(name, message string, force bool) error {
	if r.CurrentBranch.Name != r.DefaultBranch.Name {
		logging.Instance().Warnf("creating tag based on a non-default branch. it is recommended to only create tags from a base branch. current branch: %s", r.CurrentBranch.Name)
	}

	var tagCmd *exec.Cmd
	if force {
		tagCmd = exec.Command("git", "tag", "-a", name, "--force", "-m", message)
	} else {
		tagCmd = exec.Command("git", "tag", "-a", name, "-m", message)
	}

	stdout, err := tagCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v. %s", err, common.CleanstdoutMultiline(stdout))
	}

	return nil
}

func (r *Repository) String() string {
	return fmt.Sprintf("Repository: { Name: %s, VersionPrefix: '%s', Default Branch: '%s', Current Branch: '%s', Last Tag: %s }",
		r.Name,
		r.VersionPrefix,
		r.DefaultBranch,
		r.CurrentBranch,
		r.LastTag)
}