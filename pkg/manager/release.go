package manager

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type Asset struct {
	Name string `json:"name"`
	Browser_download_url string `json:"browser_download_url"`
	Size int `json:"size"`
}

type Release struct {
	Tag string `json:"tag_name"`
	Assets []Asset `json:"assets"`
}

var ArchMap = map[string]string{
    "386": "32bit",
    "amd64": "64bit",
    "arm": "ARMv6",
    "arm64": "ARM64",
}

func readData(reader io.ReadCloser) ([]byte, error) {
	content := []byte{}
	b := make([]byte, 50)
	for {
		count, readErr := reader.Read(b)
		if count > 0 {
			content = append(content, b[:count]...)
		}

		if readErr == io.EOF {
			break
		}

		if readErr != nil {
			return nil, readErr
		}
	}

	return content, nil
}

func downloadFromURL(ctx context.Context, task *TaskProgress, path, url string) {
	defer task.SetDone(path)
	
	out, err := os.Create(path)
	if err != nil {
		task.SetError(err)
		return
	}
	defer out.Close()

	res, err := http.Get(url)
	if err != nil {
		task.SetError(err)
		return
	}
	defer res.Body.Close()

	progress := TransferData(ctx, res.Body, out)
	for value := range progress {
		task.SetProgress(value)
	}
}

func GetLatestRelease() (*Release, error) {
	res, err := http.Get("https://api.github.com/repos/arduino/arduino-cli/releases/latest")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := readData(res.Body)
	if err != nil {
		return nil, err
	}

	rel := &Release{}
	jsonError := json.Unmarshal(content, rel)
	if jsonError != nil {
		return nil, jsonError
	}

	return rel, nil
}

func GetReleases() ([]*Release, error) {
	res, err := http.Get("https://api.github.com/repos/arduino/arduino-cli/releases")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := readData(res.Body)
	if err != nil {
		return nil, err
	}

	rel := []*Release{}
	jsonError := json.Unmarshal(content, &rel)
	if jsonError != nil {
		return nil, jsonError
	}

	return rel, nil
}

func DownloadRelease(ctx context.Context, path, tag string) (*TaskProgress, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("path must be a directory")
	}

	releases, err := GetReleases()
	if err != nil {
		return nil, err
	}

	// Map GOARCH to arduino-cli naming convention
	arch := ArchMap[runtime.GOARCH]

	platform := runtime.GOOS
	found := ""
	url := ""
	size := 0
	for _, rel := range releases {
		if rel.Tag == tag {
			for _, asset := range rel.Assets {
				download := strings.ToLower(asset.Browser_download_url)
				size = asset.Size
				if strings.Contains(download, arch) && strings.Contains(download, platform) {
					found = path + string(os.PathSeparator) + asset.Name
					url = asset.Browser_download_url
					break
				}
			}
		}
	}

	if found != "" {
		task := NewTaskProgress(size)
		go downloadFromURL(ctx, task, found, url)
		return task, nil
	}

	return nil, errors.New("file not found")
}
