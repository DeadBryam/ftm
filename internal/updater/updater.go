package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	githubAPI       = "https://api.github.com"
	defaultRepo     = "sthbryan/ftm"
	apiTimeout      = 15 * time.Second
	downloadTimeout = 5 * time.Minute
	userAgent       = "ftm-updater"
)

type Release struct {
	TagName    string  `json:"tag_name"`
	Prerelease bool    `json:"prerelease"`
	Draft      bool    `json:"draft"`
	HTMLURL    string  `json:"html_url"`
	Assets     []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Info struct {
	CurrentVersion string
	LatestVersion  string
	Tag            string
	AssetURL       string
	AssetName      string
	HasUpdate      bool
	ReleaseURL     string
}

type Updater struct {
	repo string
	cli  *http.Client
}

func New(repo string) *Updater {
	if repo == "" {
		repo = defaultRepo
	}
	return &Updater{
		repo: repo,
		cli:  &http.Client{Timeout: apiTimeout},
	}
}

func (u *Updater) Check(currentVersion string) (*Info, error) {
	rel, err := u.fetchLatest()
	if err != nil {
		return nil, err
	}

	latestVersion := strings.TrimPrefix(rel.TagName, "v")
	cur := strings.TrimPrefix(currentVersion, "v")

	assetName := platformAssetName()
	assetURL := ""
	for _, a := range rel.Assets {
		if a.Name == assetName {
			assetURL = a.BrowserDownloadURL
			break
		}
	}
	if assetURL == "" {
		return nil, fmt.Errorf("no asset %q in release %s (available: %s)",
			assetName, rel.TagName, listAssetNames(rel.Assets))
	}

	return &Info{
		CurrentVersion: cur,
		LatestVersion:  latestVersion,
		Tag:            rel.TagName,
		AssetURL:       assetURL,
		AssetName:      assetName,
		HasUpdate:      compareSemver(latestVersion, cur) > 0,
		ReleaseURL:     rel.HTMLURL,
	}, nil
}

func (u *Updater) fetchLatest() (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", githubAPI, u.repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := u.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("github API status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	if rel.Prerelease || rel.Draft {
		return nil, fmt.Errorf("latest release is prerelease/draft")
	}
	if rel.TagName == "" {
		return nil, fmt.Errorf("release has no tag_name")
	}
	return &rel, nil
}

func (u *Updater) Apply(info *Info) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("self-update on Windows is not supported yet; please reinstall from %s", info.ReleaseURL)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate binary: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(execPath); err == nil {
		execPath = resolved
	}

	tmp, err := os.CreateTemp(filepath.Dir(execPath), "ftm-update-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpPath := tmp.Name()

	if err := downloadFile(info.AssetURL, tmp); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("download: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("close temp: %w", err)
	}
	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("chmod: %w", err)
	}

	oldPath := execPath + ".old"
	_ = os.Remove(oldPath)
	if err := os.Rename(execPath, oldPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("backup old binary: %w", err)
	}
	if err := os.Rename(tmpPath, execPath); err != nil {
		_ = os.Rename(oldPath, execPath)
		return fmt.Errorf("install new binary: %w", err)
	}

	go func() {
		time.Sleep(2 * time.Second)
		_ = os.Remove(oldPath)
	}()

	if runtime.GOOS == "darwin" {
		_ = exec.Command("xattr", "-d", "com.apple.quarantine", execPath).Run()
	}
	return nil
}

func downloadFile(url string, dst *os.File) error {
	cli := &http.Client{Timeout: downloadTimeout}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/octet-stream")

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	_, err = io.Copy(dst, resp.Body)
	return err
}

func platformAssetName() string {
	return fmt.Sprintf("ftm-%s-%s", osAlias(), archAlias())
}

func osAlias() string {
	switch runtime.GOOS {
	case "darwin":
		return "macos"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	}
	return runtime.GOOS
}

func archAlias() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x64"
	case "arm64":
		return "arm64"
	}
	return runtime.GOARCH
}

func compareSemver(a, b string) int {
	pa := parseSemver(a)
	pb := parseSemver(b)
	for i := 0; i < 3; i++ {
		if pa[i] < pb[i] {
			return -1
		}
		if pa[i] > pb[i] {
			return 1
		}
	}
	return 0
}

func parseSemver(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var out [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		seg := strings.SplitN(parts[i], "-", 2)[0]
		seg = strings.SplitN(seg, "+", 2)[0]
		n, _ := strconv.Atoi(seg)
		out[i] = n
	}
	return out
}

func listAssetNames(assets []Asset) string {
	names := make([]string, 0, len(assets))
	for _, a := range assets {
		names = append(names, a.Name)
	}
	return strings.Join(names, ", ")
}
