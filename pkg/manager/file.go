package manager

import (
	// "archive/tar"
	"archive/zip"
	// "compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func TransferFileContents(ctx context.Context, reader io.Reader, writer io.Writer, task *TaskProgress) {
	b := make([]byte, 32768)
	total := 0
	for {
		select {
		case <- ctx.Done():
			task.SetError(errors.New("interrupted"))
			task.SetDone()
			return
		default:
			count, err := reader.Read(b)
			if count > 0 {
				_, writeErr := writer.Write(b[:count])
				if writeErr != nil {
					task.SetError(writeErr)
					task.SetDone()
					return
				}
				total += count
			}

			if err == io.EOF {
				task.SetDone()
				return
			}

			if err != nil  {
				task.SetError(err)
				task.SetDone()
				return
			}

			task.SetProgress(total)
		}
	}
}

func CompressFileZip(ctx context.Context, path string) (*TaskProgress, error) {
	pathZip := ""
	// Replace extension
	index := strings.LastIndex(path, ".")
	if index >= 0 {
		pathZip = path[:index+1] + "zip"
	} else {
		pathZip = path + ".zip"
	}

	// Output file
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	task := NewTaskProgress(int(info.Size()))
	go func () {
		fileIn, inErr := os.Open(path)
		if inErr != nil {
			task.SetError(inErr)
			task.SetDone()
			return
		}
		defer fileIn.Close()

		fileOut, outErr := os.Create(pathZip)
		if outErr != nil {
			task.SetError(outErr)
			task.SetDone()
			return
		}
		defer fileOut.Close()

		header, headerErr := zip.FileInfoHeader(info)
		if headerErr != nil {
			task.SetError(headerErr)
			task.SetDone()
			return
		}

		header.Name = filepath.Base(path)
		header.Method = zip.Deflate
	
		writer := zip.NewWriter(fileOut)
		// f, err := writer.Create(path)
		f, err := writer.CreateHeader(header)
		if err != nil {
			task.SetError(err)
			task.SetDone()
			return
		}
		defer writer.Close()
	
		TransferFileContents(ctx, fileIn, f, task)
	}()

	return task, nil
}

func CompressFileTarGz(ctx context.Context, path string) error {
	return nil
}

func DecompressFile(ctx context.Context, task *TaskProgress, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if strings.HasSuffix(path, ".zip") {
		fmt.Println("ZIP")
	} else if strings.HasSuffix(path, "tar.gz") {
		fmt.Println("TAR.GZ")
	}

	return nil
}