package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"git.gmantaos.com/haath/Goirate/pkg/utils"
	update "github.com/inconshreveable/go-update"
)

// LatestReleaseURL is the API endpoint through which information about the latest release is retrieved.
const LatestReleaseURL = "https://api.github.com/repos/gmantaos/Goirate/releases/latest"

// UpdateCommand defines the update command and holds its options.
type UpdateCommand struct {
}

type releaseResponse struct {
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

// Execute is the callback of the update command.
func (cmd *UpdateCommand) Execute(args []string) error {

	currentVersion, valid := parseVersion(VERSION)

	if !valid && os.Getenv("GOIRATE_DEBUG") != "true" {

		return fmt.Errorf("cannot automatically update from current version: %v", VERSION)
	}

	var resp releaseResponse

	err := utils.HTTPGetJSON(LatestReleaseURL, &resp)

	if err != nil {
		return err
	}

	latestVersion, valid := parseVersion(resp.TagName)

	if !valid {

		return fmt.Errorf("there appears to be an error with the GitHub releases for the project")
	}

	if os.Getenv("GOIRATE_DEBUG") == "true" {

		log.Printf("%s -> %s [%s.%s]", currentVersion, latestVersion, runtime.GOOS, runtime.GOARCH)
	}

	if latestVersion.moreRecentThan(currentVersion) {

		log.Printf("Updating to version: %v\n", latestVersion)

		binaryURL, exists := resp.getDownloadURLForArch()

		if !exists {

			return fmt.Errorf("architecture missing from GitHub releases")
		}

		return doUpdate(binaryURL)

	}

	log.Printf("Up to date. (%v)", currentVersion)

	return nil
}

func doUpdate(binaryURL string) error {

	resp, err := http.Get(binaryURL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = update.Apply(resp.Body, update.Options{})

	return err
}

func parseVersion(versionString string) (version, bool) {

	r, _ := regexp.Compile(`v?(\d+).(\d+).(\d+)`)
	m := r.FindStringSubmatch(strings.TrimSpace(versionString))

	var major, minor, patch uint64
	valid := true

	if len(m) > 0 {

		var err error

		major, err = strconv.ParseUint(m[1], 10, 64)

		if err != nil {
			valid = false
		}

		minor, err = strconv.ParseUint(m[2], 10, 64)

		if err != nil {
			valid = false
		}

		patch, err = strconv.ParseUint(m[3], 10, 64)

		if err != nil {
			valid = false
		}
	}

	return version{major, minor, patch}, valid
}

func (ver version) moreRecentThan(other version) bool {

	return ver.Major > other.Major ||
		(ver.Major == other.Major && ver.Minor > other.Minor) ||
		(ver.Major == other.Major && ver.Minor == other.Minor && ver.Patch > other.Patch)
}

func (ver version) String() string {

	return fmt.Sprintf("%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
}

func (resp *releaseResponse) getDownloadURLForArch() (string, bool) {

	os := runtime.GOOS
	arch := runtime.GOARCH

	for _, asset := range resp.Assets {

		if strings.Contains(asset.Name, os) && strings.Contains(asset.Name, arch) {

			return asset.DownloadURL, true
		}
	}

	return "", false
}
