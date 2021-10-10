package manager

import (
	// "archive/tar"
	"archive/zip"
	// "compress/gzip"
	"context"
	"errors"
	"io"
	"os"
)

func GetUncompressedZipSize(files []*zip.File) int {
	total := 0
	for _, file := range files {
		total += int(file.FileInfo().Size())
	}
	return total
}

func GetFilesTotalSize(paths []string) (int, error) {
	total := 0
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return total, err
		}

		total += int(info.Size())
	}
	return total, nil
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
					return
				}
			}
		}
	}()

	return progress
}

func CompressFileZip(ctx context.Context, name string, paths []string) (*TaskProgress, error) {
	size, err := GetFilesTotalSize(paths)
	if err != nil {
		return nil, err
	}

	fileOut, outErr := os.Create(name)
	if outErr != nil {
		return nil, outErr
	}
	
	task := NewTaskProgress(size)
	go func () {
		defer fileOut.Close()
		defer task.SetDone()

		writer := zip.NewWriter(fileOut)
		defer writer.Close()

		totalSize := 0
		for _, path := range paths {
			fileIn, inErr := os.Open(path)
			if inErr != nil {
				task.SetError(inErr)
				return
			}

			f, err := writer.Create(path)
			if err != nil {
				task.SetError(err)
				return
			}
		
			count := 0
			progress := TransferData(ctx, fileIn, f)
			for value := range progress {
				count = value
			}
			totalSize += count
			task.SetProgress(totalSize)

			fileIn.Close()
		}
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
		
		if task.GetProgress() != task.GetTotal() {
			task.SetError(errors.New("uncompressed size does not match"))
		}
	}()

	return task, nil
}