package jorge

import (
	"bufio"
	"fmt"
	"os"
)

func Contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func ExistsInFile(filePath string, element string) (bool, error) {
	readFile, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	found := false
	for fileScanner.Scan() {
		if fileScanner.Text() == element {
			found = true
		}
	}

	return found, nil
}

func AppendToFile(filePath string, element string) error {

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err2 := file.WriteString(fmt.Sprintf("\n%s\n", element))

	if err2 != nil {
		return err2

	} else {
		return nil
	}
}
