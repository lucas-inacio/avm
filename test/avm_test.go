package manager

import (
	"context"
	"testing"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func TestDownloadRelease(t *testing.T) {
	
}

func TestCompressFileZip(t *testing.T) {
	task, err := manager.CompressFileZip(context.Background(), "agoravai.zip", []string{"tmp/arduino-cli.exe", "tmp/LICENSE.txt"})
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
}

func TestDecompressFileZip(t *testing.T) {
	task, err := manager.DecompressFileZip(context.Background(), "../arduino-cli_0.19.2_Windows_64bit.zip")
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
}