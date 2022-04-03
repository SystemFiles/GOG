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
	"sykesdev.ca/gog/lib/semver"
)

var ctx = context.Background()
var Version string = "1.0.2"

type Updater struct {
	client *github.Client

	currentVersion semver.Semver
	updateVersion semver.Semver

	repoOwner string
	repoName string

	latestRelease *github.RepositoryRelease

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
		repoOwner: "systemfiles",
		repoName: "gog",
		binaryOs: runtime.GOOS,
		binaryArch: runtime.GOARCH,
		binaryLocation: binaryPath,
	}

	u.currentVersion = semver.MustParse(Version)
	if tag == "" {
		if err := u.getLatestVersionAndRelease(); err != nil {
			return nil, err
		}
	} else {
		u.updateVersion = semver.MustParse(tag)
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

func (u *Updater) getLatestVersionAndRelease() error {
	releases, _, err := u.client.Repositories.ListReleases(ctx, u.repoOwner, u.repoName, &github.ListOptions{})
	if err != nil {
		return errors.New("failed to get latest GOG version. " + err.Error())
	}

	u.updateVersion = semver.MustParse(*releases[0].TagName)
	u.latestRelease = releases[0]

	return nil
}

func (u *Updater) getLatestReleaseAsset() (*github.ReleaseAsset, error) {
	if u.latestRelease == nil {
		return nil, errors.New("cannot get latest release asset since 'latestRelease' is not defined")
	}

	for _, asset := range u.latestRelease.Assets {
		extension := "tar.gz"
		if u.binaryOs == "windows" {
			extension = "zip"
		}
		if strings.Contains(*asset.Name, fmt.Sprintf("%s-%s-%s-%s.%s", u.repoName, u.updateVersion, u.binaryOs, u.binaryArch, extension)) {
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

	req, err := http.NewRequestWithContext(ctx, "GET", downloadUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")

	return req.GetBody()
}

func (u *Updater) Update() error {

	tarFile, err := u.downloadLatestReleaseBinary()
	if err != nil {
		return err
	}

	tData, err := UntarBinary(tarFile, u.repoName)
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
		return nil
	}

	return nil
}

func (u *Updater) String() string {
	return fmt.Sprintf("Updater<%s - %s>(%s->%s)", u.repoName, u.repoOwner, u.currentVersion, u.updateVersion)
}







