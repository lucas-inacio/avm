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
	"time"
	"unicode"
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

func downloadFromURL(ctx context.Context, task *TaskProgress, path, url string, size int) {
	out, err := os.Create(path)
	if err != nil {
		task.SetError(err)
		task.SetDone()
		return
	}
	defer out.Close()

	res, err := http.Get(url)
	if err != nil {
		task.SetError(err)
		task.SetDone()
		return
	}
	defer res.Body.Close()

	total := 0
	b := make([]byte, 32768)

	ticker := time.NewTicker(time.Microsecond)
	for {
		select {
		case <- ticker.C:
			readCount, readErr := res.Body.Read(b)
			if readCount > 0 {
				for readCount > 0 {
					writeCount, writeErr := out.Write(b[:readCount])
					if writeErr != nil {
						task.SetError(writeErr)
						return
					}
					readCount -= writeCount
					total += writeCount
				}
			}
	
			if readErr == io.EOF {
				task.SetDone()
				return
			}
			
			if readErr != nil {
				task.SetError(readErr)
				task.SetDone()
				return
			}

			task.SetProgress(float32(total) / float32(size))
		case <- ctx.Done():
			task.SetError(errors.New("interrupted"))
			task.SetDone()
			return
		}
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
	releases, err := GetReleases()
	if err != nil {
		return nil, err
	}

	// Extract numbers from GOARCH string
	arch := func () string {
		archString := ""
		for _, character := range runtime.GOARCH {
			if unicode.IsDigit(character) {
				archString += string(character)
			}
		}
		return archString
	}()

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
		task := NewTaskProgress()
		go downloadFromURL(ctx, task, found, url, size)
		return task, nil
	}

	return nil, errors.New("file not found")
}