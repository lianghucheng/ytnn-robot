package common

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func ReadFile(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	defer f.Close()
	var names []string
	if err != nil {
		return names, err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err == nil {
			name := strings.TrimSpace(line)
			if name != "" {
				names = append(names, name)
			}
		} else {
			if err == io.EOF {
				return names, nil
			}
			return names, err
		}

	}
	return names, nil
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	//log.Debug("f: %v", filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func WriteMapToFile(m map[string]int, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("create file error: %v\n", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for k := range m {
		fmt.Fprintln(w, k)
	}
	return w.Flush()
}
