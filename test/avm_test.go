package manager

import (
	"context"
	"testing"

	"github.com/lucas-inacio/avm/pkg/manager"
)

func TestDownloadRelease(t *testing.T) {
	
}

func TestCompressFileZip(t *testing.T) {
	task, err := manager.CompressFileZip(context.Background(), "data/test_file.txt")
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

func TestDecompress(t *testing.T) {

}