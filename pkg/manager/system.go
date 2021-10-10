package manager

import (
	"errors"
	"io"
	"os/exec"
	"regexp"
)

func GetArduinoVersion() (string, error) {
	_, err := exec.LookPath("arduino-cli")
	if err != nil {
		return "", err
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