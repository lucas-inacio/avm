package manager

import (
	"context"
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
	"unicode"

	"github.com/lucas-inacio/avm/pkg/manager"
)

const (
	DumbFileName1    = "test_file_1.txt"
	DumbFileContent1 = "This is a test string"
	DumbFileName2    = "test_file_2.txt"
	DumbFileContent2 = "This is another test string"
	OutputFileName   = "output.zip"
	TmpDirectory     = "tmp"
)

func createDumbFiles() error {
	file, err := os.Create(TmpDirectory + string(os.PathSeparator) + DumbFileName1)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write([]byte(DumbFileContent1))

	file2, err2 := os.Create(TmpDirectory + string(os.PathSeparator) + DumbFileName2)
	if err2 != nil {
		return err2
	}
	defer file2.Close()
	file2.Write([]byte(DumbFileContent2))

	return nil
}

func removeDumbFiles() error {
	err := os.Remove(TmpDirectory + string(os.PathSeparator) + DumbFileName1)
	if err != nil {
		return err
	}

	err = os.Remove(TmpDirectory + string(os.PathSeparator) + DumbFileName2)
	if err != nil {
		return err
	}

	return nil
}

func checkFileContent(path, content string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	data := ""
	b := make([]byte, info.Size())
	for {
		count, readErr := file.Read(b)
		if count > 0 {
			data += string(b[:count])
		}

		if readErr == io.EOF {
			break
		}

		if readErr != nil {
			return readErr
		}
	}

	if data != content {
		return errors.New("file contents do not match")
	}

	return nil
}

func TestDownloadRelease(t *testing.T) {
	rel, err := manager.GetLatestRelease()
	if err != nil {
		t.Error(err)
	}

	if err := os.Mkdir(TmpDirectory, os.ModeDir); err != nil {
		t.Error(err)
	}

	task, err := manager.DownloadRelease(context.Background(), TmpDirectory, rel.Tag)
	if err != nil {
		t.Error(err)
	}

	<- task.Done()

	dirs, err := os.ReadDir(TmpDirectory)
	if err != nil {
		t.Error(err)
	}

	if len(dirs) != 1 {
		t.Error(errors.New("invalid directory tree; expected " + TmpDirectory + " to have only one file"))
	}

	name := ""
	found := false
	for _, item := range dirs {

		name = func () string {
			result := ""
			for _, r := range item.Name() {
				result += string(unicode.ToLower(r))
			}
			return result
		}()
		
		if strings.Contains(name, runtime.GOOS) {
			found = true
		}
	}

	if !found {
		t.Error(errors.New("file was not downloaded"))
	}

	if err := os.RemoveAll(TmpDirectory); err != nil {
		t.Error(err)
	}
}

func TestDecompressFileZip(t *testing.T) {
	task, err := manager.DecompressFileZip(
		context.Background(), TmpDirectory + string(os.PathSeparator) + OutputFileName)
	if err != nil {
		t.Error(err)
	}
	
	<- task.Done()
	
	if task.GetError() != nil {
		t.Error(task.GetError())
	}
	
	if task.GetProgress() != task.GetTotal() {
		t.Fail()
	}

	check1 := checkFileContent(
		TmpDirectory + string(os.PathSeparator) + DumbFileName1,
		DumbFileContent1,
	)
		
	if check1 != nil {
		t.Error(check1)
	}

	check2 := checkFileContent(
		TmpDirectory + string(os.PathSeparator) + DumbFileName2,
		DumbFileContent2,
	)
	if check2 != nil {
		t.Error(check2)
	}

	if err := os.RemoveAll("tmp"); err != nil {
		t.Error(err)
	}
}