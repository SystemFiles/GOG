package command

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"sykesdev.ca/gog/lib"
)

func usage() {
	lib.GetLogger().Info("Usage: gog feature <jira_name> [comment]")
}

func GitIsValidRepo() bool {
	cmd := exec.Command("git", "status")
	_, err := cmd.Output()
	
	return err == nil
}

func GitBranchExists(branch string) bool {
	cmd := exec.Command("git", "branch", "|", "egrep", branch)
	_, err := cmd.Output()
	
	return err == nil
}

func GitUnstagedCommits() bool {
	cmd := exec.Command("git", "diff-index", "--quiet", "HEAD")
	_, err := cmd.Output()

	return err != nil
}

func GitGetCurrentBranch() (string, error) {
	cmd := exec.Command("bash", "-c", "git branch | grep '*' | cut -d' ' -f2")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

func GitPullChanges() error {
	cmd := exec.Command("git", "pull")
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func GitCreateBranch(name string, checkout bool) error {
	if checkout {
		cmd := exec.Command("git", "checkout", "-b", name)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command("git", "branch", name)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}

	return nil
}

func GitOriginDefaultBranch() (string, error) {
	defaultBranchCmd := exec.Command("bash", "-c", "git remote show origin | sed -n '/HEAD branch/s/.*: //p'")
	defaultBranch, err := defaultBranchCmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(defaultBranch)), nil
}

func GitCheckoutDefaultBranch() error {
	defaultBranch, err := GitOriginDefaultBranch()
	if err != nil {
		return err
	}

	checkoutCmd := exec.Command("git", "checkout", defaultBranch)
	_, err = checkoutCmd.Output()
	if err != nil {
		return err
	}

	if err = GitPullChanges(); err != nil {
		return err
	}

	return nil
}

func GOGNewFeature(cwd string, feature *lib.Feature) error {
	if !PathExists(cwd + "/.gog") {
		if err := os.MkdirAll(cwd + "/.gog", 0700); err != nil {
			return err
		}
	}

	f, err := os.Create(cwd + "/.gog/feature.json")
	if err != nil {
		return err
	}
	defer f.Close()

	featureBytes, err := json.Marshal(*feature)
	if err != nil {
		return err
	}

	_, err = f.Write(featureBytes)
	if err != nil {
		return err
	}

	return nil
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CleanFeature(cwd string, feature *lib.Feature) error {
	if err := GitCheckoutDefaultBranch(); err != nil {
		return err
	}

	cmd := exec.Command("git", "branch", "-D", feature.Jira)
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	os.RemoveAll(cwd + "/.gog/")

	return nil
}

func ExecFeature(opts []string) {
	if len(opts) < 1 {
		usage()
		os.Exit(0)
	}

	jira := opts[0]
	comment := strings.Join(opts[1:], " ")
	if comment == "" {
		comment = "Feature Branch"
	}

	workingDir, err := os.Getwd()
	if err != nil {
		lib.GetLogger().Fatal("Failed to get working directory from path")
	}
	GOGDir := fmt.Sprintf("%s/%s", workingDir, ".gog")
	feature, err := lib.NewFeature(jira, comment)
	if err != nil {
		lib.GetLogger().Error("Failed to create feature")
		lib.GetLogger().Error(fmt.Sprintf("Reason: %v", err))
		os.Exit(1)
	}

	if !GitIsValidRepo() {
		lib.GetLogger().Error("The current directory does not contain a valid git repo ...")
		os.Exit(1)
	}

	initial_branch, err := GitGetCurrentBranch()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to get current branch of git repository. %v", err))
		os.Exit(1)
	}

	if GitBranchExists(feature.Jira) {
		lib.GetLogger().Error(fmt.Sprintf("There is already a branch in this repo titled %s", feature.Jira))
		os.Exit(1)
	}

	if PathExists(GOGDir) {
		lib.GetLogger().Error(fmt.Sprintf("%s already exists ... there could already be a feature here. Please fix this and try again.", GOGDir))
		os.Exit(1)
	}

	if GitUnstagedCommits() {
		lib.GetLogger().Error(fmt.Sprintf("There is unstaged commits on your current branch (%s). For your safety, please stage or discard the changes to continue. %v", initial_branch, err))
		os.Exit(1)
	}

	if !lib.StringInSlice(opts, "from-feature") {
		if err := GitCheckoutDefaultBranch(); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to checkout default branch for repo. %v", err))
			os.Exit(1)
		}
	}

	if err := GitCreateBranch(feature.Jira, true); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to create or checkout new feature branch, %s. %v", feature.Jira, err))
		os.Exit(1)
	}

	if err := GOGNewFeature(workingDir, feature); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to create feature tracking file (%v) ... will exit cleanly", err))
		if err := CleanFeature(workingDir, feature); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to exit cleanly ... %v", err))
		}
		os.Exit(1)
	}

	lib.GetLogger().Info(fmt.Sprintf("Successfully created feature %s!", jira))
}