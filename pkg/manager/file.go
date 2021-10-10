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

func GetUncompressedZipSize(files []*zip.File) int {
	total := 0
	for _, file := range files {
		total += int(file.FileInfo().Size())
		fmt.Println(file.FileInfo().Size())
	}
	return total
}

func TransferData(ctx context.Context, reader io.Reader, writer io.Writer) chan int {
	progress := make(chan int)
	go func() {
		defer close(progress)

		b := make([]byte, 32768)
		total := 0
		for {
			select {
			case <- ctx.Done():
				return
			default:
				count, err := reader.Read(b)
				if count > 0 {
					_, writeErr := writer.Write(b[:count])
					if writeErr != nil {
						return
					}
					total += count
				}
				
				progress <- total

				if err == io.EOF {
					return
				}
	
				if err != nil  {
					fmt.Println("Error")
					return
				}
	
			}
		}
	}()

	return progress
}

func TransferFileContents(ctx context.Context, reader io.Reader, writer io.Writer, task *TaskProgress) {
	b := make([]byte, 32768)
	total := 0
	defer task.SetDone()
	for {
		select {
		case <- ctx.Done():
			task.SetError(errors.New("interrupted"))
			return
		default:
			count, err := reader.Read(b)
			if count > 0 {
				_, writeErr := writer.Write(b[:count])
				if writeErr != nil {
					task.SetError(writeErr)
					return
				}
				total += count
			}

			if err == io.EOF {
				return
			}

			if err != nil  {
				task.SetError(err)
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

func DecompressFileZip(ctx context.Context, path string) (*TaskProgress, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	size := GetUncompressedZipSize(reader.File)
	fmt.Println("Size:", size)

	totalSize := 0
	count := 0
	task := NewTaskProgress(size)
	go func () {
		defer reader.Close()
		defer task.SetDone()

		for _, file := range reader.File {
			in, inErr := file.Open()
			if inErr != nil {
				task.SetError(inErr)
				return
			}
			defer in.Close()

			out, outErr := os.Create(file.Name)
			if outErr != nil {
				task.SetError(outErr)
				return
			}
			defer out.Close()

			progress := TransferData(context.Background(), in, out)
			for value := range progress {
				count = value
			}
			totalSize += count
			task.SetProgress(totalSize)
		}
		fmt.Println("After:", totalSize)
	}()

	return task, nil
}