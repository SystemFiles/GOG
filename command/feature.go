package command

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

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

	return err == nil
}

func GitGetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "|", "grep", "'*'", "|", "cut", "-d' '", "-f2")
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

func GitCheckoutDefaultBranch() error {
	default_branch_cmd := exec.Command("git", "remote", "show", "origin", "|", "sed", "-n", "'/HEAD branch/s/.*: //p'")
	default_branch, err := default_branch_cmd.Output()
	if err != nil {
		return err
	}

	checkout_cmd := exec.Command("git", "checkout", string(default_branch))
	_, err = checkout_cmd.Output()
	if err != nil {
		return err
	}

	// 7. Pull any changes from origin on the default branch
	if err = GitPullChanges(); err != nil {
		return err
	}

	return nil
}

func GOGNewFeature(cwd string, feature *lib.Feature) error {
	if !PathExists(cwd) {
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

func ExecFeature(opts []string) {
	if len(opts) < 1 {
		usage()
		os.Exit(0)
	}

	jira := opts[0]
	comment := opts[1]
	if comment == "" {
		comment = "Feature Branch"
	}

	working_dir, err := os.Getwd()
	if err != nil {
		lib.GetLogger().Fatal("Failed to get working directory from path")
	}
	gog_dir := fmt.Sprintf("%s/%s", working_dir, ".gog")

	feature, err := lib.NewFeature(jira, comment)
	if err != nil {
		lib.GetLogger().Error("Failed to create feature")
		lib.GetLogger().Error(fmt.Sprintf("Reason: %v", err))
		os.Exit(1)
	}

	// 2. Check if inside git repository
	if !GitIsValidRepo() {
		lib.GetLogger().Error("The current directory does not contain a valid git repo ...")
		os.Exit(1)
	}

	initial_branch, err := GitGetCurrentBranch()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to get current branch of git repository. %v", err))
		os.Exit(1)
	}

	// 3. Check if there is already a branch with the same name as the one to create
	if GitBranchExists(feature.Jira) {
		lib.GetLogger().Error(fmt.Sprintf("There is already a branch in this repo titled %s", feature.Jira))
		os.Exit(1)
	}

	// 4. Check if there is an existing .gog folder
	if PathExists(gog_dir) {
		lib.GetLogger().Error(fmt.Sprintf("%s already exists ... there could already be a feature here. Please fix this and try again.", gog_dir))
		os.Exit(1)
	}

	// 5. Check if there are any unstaged commits
	if GitUnstagedCommits() {
		lib.GetLogger().Error(fmt.Sprintf("There is unstaged commits on your current branch (%s). For your safety, please stage or discard the changes to continue.", initial_branch))
		os.Exit(1)
	}

	// 6. Check for from-feature condition
	if !lib.StringInSlice(opts, "from-feature") {
		if err := GitCheckoutDefaultBranch(); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to checkout default branch for repo. %v", err))
			os.Exit(1)
		}
	}

	// 8. Create a branch using the feature object (Jira parameter as the branch name)
	// 9. Checkout to the new branch
	if err := GitCreateBranch(feature.Jira, true); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to create or checkout new feature branch, %s. %v", feature.Jira, err))
	}

	// 10. Create a .gog/feature.json file
	GOGNewFeature(working_dir, feature)

	lib.GetLogger().Info(fmt.Sprintf("Successfully created feature %s!", jira))
}