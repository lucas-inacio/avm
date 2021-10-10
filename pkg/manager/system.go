package manager

import (
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
)

type ArduinoNotFoundError struct {
	message string
}

func NewArduinoNotFoundError() *ArduinoNotFoundError {
	return &ArduinoNotFoundError{
		"arduino-cli not found",
	}
}

func (err *ArduinoNotFoundError) Error() string {
	return err.message
}

func GetArduinoVersion() (string, error) {
	_, err := GetArduinoDir()
	if err != nil {
		return "", NewArduinoNotFoundError()
	}

	cmd := exec.Command("arduino-cli", "version")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	output := []byte{}
	b := make([]byte, 128)
	for {
		count, inErr := stdout.Read(b)
		if count > 0 {
			output = append(output, b[:count]...)
		}

		if inErr == io.EOF {
			break
		}

		if inErr != nil {
			return "", inErr
		}
	}

	reg := regexp.MustCompile(`\d+\.\d+\.\d+(-\w+)?`)
	if data := reg.Find(output); data != nil {
		return string(data), nil
	}

	return "", errors.New("could not determine arduino-cli version")
}

func GetArduinoDir() (string, error) {
	dir, err := exec.LookPath("arduino-cli")
	if err != nil {
		return "", NewArduinoNotFoundError()
	}

	return filepath.Dir(dir), nil
}