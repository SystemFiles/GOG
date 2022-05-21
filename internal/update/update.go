package update

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/google/go-github/v43/github"
	"sykesdev.ca/gog/internal/semver"
)

var ctx = context.Background()
var Version string

type Updater struct {
	client *github.Client

	currentVersion semver.Semver
	updateVersion semver.Semver

	repoOwner string
	repoName string

	updateRelease *github.RepositoryRelease

	binaryOs string
	binaryArch string
	binaryLocation string
}

func NewUpdater(tag string) (*Updater, error) {
	binaryPath, err := os.Executable()
	if err != nil {
		panic("failed to get executable path")
	}

	u := &Updater{
		client: github.NewClient(nil),
		repoOwner: "SystemFiles",
		repoName: "GOG",
		binaryOs: runtime.GOOS,
		binaryArch: runtime.GOARCH,
		binaryLocation: binaryPath,
	}

	u.currentVersion = semver.MustParse(Version)
	if tag == "" {
		u.updateVersion, err = u.getLatestVersion()
		if err != nil {
			return nil, err
		}
	} else {
		u.updateVersion, err = semver.Parse(tag)
		if err != nil {
			return nil, err
		}
	}

	if err := u.getReleaseForVersion(u.updateVersion); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *Updater) Client() *github.Client {
	return u.client
}

func (u *Updater) RepoOwner() string {
	return u.repoOwner
}

func (u *Updater) RepoName() string {
	return u.repoName
}

func (u *Updater) BinaryOS() string {
	return u.binaryOs
}

func (u *Updater) BinaryArch() string {
	return u.binaryArch
}

func (u *Updater) BinaryLocation() string {
	return u.binaryLocation
}

func (u *Updater) CurrentVersion() semver.Semver {
	return u.currentVersion
}

func (u *Updater) UpdateVersion() semver.Semver {
	return u.updateVersion
}

func (u *Updater) getLatestVersion() (semver.Semver, error) {
	releases, _, err := u.client.Repositories.ListReleases(ctx, u.repoOwner, u.repoName, &github.ListOptions{})
	if err != nil {
		return semver.Semver{}, errors.New("failed to list project releases. " + err.Error())
	}

	return semver.MustParse(*releases[0].TagName), nil
}

func (u *Updater) getReleaseForVersion(version semver.Semver) error {
	releases, _, err := u.client.Repositories.ListReleases(ctx, u.repoOwner, u.repoName, &github.ListOptions{})
	if err != nil {
		return errors.New("failed to list project releases. " + err.Error())
	}

	for _, r := range releases {
		if *r.TagName == version.NoPrefix() {
			u.updateRelease = r
			break
		}
	}

	if u.updateRelease == nil {
		return fmt.Errorf("failed to locate the specified version (%s) in project releases", version)
	}

	return nil
}

func (u *Updater) getLatestReleaseAsset() (*github.ReleaseAsset, error) {
	if u.updateRelease == nil {
		return nil, errors.New("cannot get latest release asset since 'latestRelease' is not defined")
	}

	for _, asset := range u.updateRelease.Assets {
		extension := "tar.gz"
		if u.binaryOs == "windows" {
			extension = "zip"
		}
		if strings.Contains(*asset.Name, fmt.Sprintf("%s-%s-%s-%s.%s", u.repoName, u.updateVersion.NoPrefix(), u.binaryOs, u.binaryArch, extension)) {
			return asset, nil
		}
	}

	return nil, errors.New("latest release asset was not found for this OS")
}

func (u *Updater) downloadLatestReleaseBinary() (io.ReadCloser, error) {
	asset, err := u.getLatestReleaseAsset()
	if err != nil {
		return nil, err
	}
	downloadUrl := asset.GetBrowserDownloadURL()

	resp, err := http.Get(downloadUrl)
  if err != nil {
    return nil, err
  }

	return resp.Body, nil
}

func (u *Updater) Update() error {
	if u.currentVersion == u.updateVersion {
		return errors.New("GOG is already at the latest version")
	}

	tarFile, err := u.downloadLatestReleaseBinary()
	if err != nil {
		return err
	}
	defer tarFile.Close()

	tData, err := UntarBinary(tarFile, strings.ToLower(u.repoName))
	if err != nil {
		return err
	}

	// remove current go binary backup (if exists)
	os.Remove(u.binaryLocation + ".old")

	// rename current go binary to keep as backup
	os.Rename(u.binaryLocation, u.binaryLocation + ".old")

	// copy updated binary to replace existing go binary
	out, err := os.Create(u.binaryLocation)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, tData)
	if err != nil {
		return err
	}

	err = os.Chmod(u.binaryLocation, 0751)
	if err != nil {
		return err
	}

	return nil
}

func (u *Updater) String() string {
	return fmt.Sprintf("Updater<%s - %s>(%s->%s)", u.repoName, u.repoOwner, u.currentVersion, u.updateVersion)
}







