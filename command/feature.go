package command

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"sykesdev.ca/gog/lib"
)

func featureUsage() {
	lib.GetLogger().Info("Usage: gog feature <jira_name> <comment>")
}

func GitIsValidRepo() bool {
	cmd := exec.Command("git", "status")
	_, err := cmd.Output()
	
	return err == nil
}

func GitBranchExists(branch string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git branch | egrep %s", branch))
	_, err := cmd.Output()

	return err == nil
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

	err := feature.Save()
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

func ExecFeature() {
	if len(flag.Args()) < 3 {
		lib.GetLogger().Error("Invalid usage of gogfeature ...")
		featureUsage()
		os.Exit(2)
	}

	jira := flag.Arg(1)
	comment := strings.Join(flag.Args()[2:], " ")
	fromFeature := *flag.Bool("from-feature", false, "specifies if this feature will be based on the a current feature branch")

	if comment == "" {
		comment = "Feature Branch"
	}

	validJiraFormat, err := regexp.Match(`^[A-Z].*\-[0-9].*$`, []byte(jira))
	if err != nil {
		lib.GetLogger().Fatal(fmt.Sprintf("Failed to parse regular expression for Jira format. %v", err))
	}

	if !validJiraFormat {
		lib.GetLogger().Error("Invalid Jira format ... example of a valid format includes JIRA-0023")
		os.Exit(1)
	}

	workingDir, GOGDir := lib.GetWorkspacePaths()

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
		lib.GetLogger().Error(fmt.Sprintf("There is already a branch in this repo named %s", feature.Jira))
		os.Exit(1)
	}

	if lib.GitHasUnstagedCommits() {
		lib.GetLogger().Error(fmt.Sprintf("There is unstaged commits on your current branch (%s). For your safety, please stage or discard the changes to continue. %v", initial_branch, err))
		os.Exit(1)
	}

	if lib.GitHasUncommittedChanges() {
		lib.GetLogger().Error(fmt.Sprintf("There are staged commits on your current branch (%s) which have not been committed. %v", initial_branch, err))
		os.Exit(1)
	}

	if !fromFeature {
		if err := GitCheckoutDefaultBranch(); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to checkout default branch for repo. %v", err))
			os.Exit(1)
		}
	}

	if PathExists(GOGDir) {
		lib.GetLogger().Error(fmt.Sprintf("%s already exists ... there could already be a feature here. Please fix this and try again.", GOGDir))
		os.Exit(1)
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