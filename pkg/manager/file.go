package manager

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
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
			case <-ctx.Done():
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

				if err != nil {
					return
				}
			}
		}
	}()

	return progress
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
	go func() {
		defer reader.Close()
		defer task.SetDone(nil)

		for _, file := range reader.File {
			in, inErr := file.Open()
			if inErr != nil {
				task.SetError(inErr)
				return
			}
			defer in.Close()

			out, outErr := os.Create(filepath.Dir(path) + string(os.PathSeparator) + file.Name)
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

func DecompressFileTargz(ctx context.Context, path string) (*TaskProgress, error) {
	stream, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(stream)
	if err != nil {
		return nil, err
	}

	task := NewTaskProgress(100)
	tarReader := tar.NewReader(gzipReader)
	baseDir := filepath.Dir(path)
	go func() {
		defer stream.Close()
		defer gzipReader.Close()
		defer task.SetDone(nil)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				task.SetError(err)
				return
			}

			if header.FileInfo().IsDir() {
				if err := os.Mkdir(baseDir+string(os.PathSeparator)+header.Name, header.FileInfo().Mode()); err != nil {
					task.SetError(err)
					return
				}
				continue
			}

			file, err := os.OpenFile(
				baseDir+string(os.PathSeparator)+header.Name,
				os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
				header.FileInfo().Mode(),
			)
			if err != nil {
				task.SetError(err)
				return
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				task.SetError(err)
				return
			}
			file.Close()
		}
	}()

	return task, nil
}
