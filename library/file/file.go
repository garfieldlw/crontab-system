package file

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const BufferSize = 1000000

func ReplaceBig(inputFile, outputFile string, replace map[string]string) error {
	in, errInput := os.OpenFile(inputFile, os.O_RDONLY, 0)
	if errInput != nil {
		return errInput
	}
	defer in.Close()

	out, errOutput := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0766)
	if errOutput != nil {
		return errOutput
	}
	defer out.Close()

	br := bufio.NewReader(in)
	index := 1
	for {
		line, _, errRead := br.ReadLine()
		if errRead == io.EOF {
			break
		}
		if errRead != nil {
			return errRead
		}

		newLine := string(line[:])
		for key, value := range replace {
			newLine = strings.Replace(newLine, key, value, -1)
		}
		_, errWrite := out.WriteString(newLine + "\n")
		if errWrite != nil {
			return errWrite
		}
		index++
	}
	return nil
}

func Replace(inputFile, outputFile string, replace map[string]string) error {
	in, errInput := ioutil.ReadFile(inputFile)
	if errInput != nil {
		return errInput
	}

	newLine := string(in[:])
	for key, value := range replace {
		newLine = strings.Replace(newLine, key, value, -1)
	}

	errOutput := ioutil.WriteFile(outputFile, []byte(newLine[:]), 0766)
	if errOutput != nil {
		return errOutput
	}

	return nil
}

func CopyBig(inputFile, outputFile string) error {
	source, errSource := os.Open(inputFile)
	if errSource != nil {
		return errSource
	}
	defer source.Close()

	destination, errDestination := os.Create(outputFile)
	if errDestination != nil {
		return errDestination
	}
	defer destination.Close()

	buf := make([]byte, BufferSize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

func Copy(inputFile, outputFile string) error {
	input, errInput := ioutil.ReadFile(inputFile)
	if errInput != nil {
		return errInput
	}

	errOutput := ioutil.WriteFile(outputFile, input, 0766)
	if errOutput != nil {
		return errOutput
	}

	return nil
}

func CopyShell(input, output string, isDir bool) error {
	if runtime.GOOS == "windows" {
		command := ""
		if isDir {
			command = fmt.Sprintf("xcopy /y %v\\* %v", input, output)
		} else {
			command = fmt.Sprintf("echo f | xcopy /y %v %v", input, output)
		}

		fmt.Println(command)

		cmd := exec.Command("cmd", "/C", command)
		errStart := cmd.Start()
		if errStart != nil {
			return errStart
		}

		errWait := cmd.Wait()
		if errWait != nil {
			return errWait
		}
	} else {
		command := ""
		if isDir {
			command = fmt.Sprintf("cp -rf %v/* %v", input, output)
		} else {
			command = fmt.Sprintf("cp -f %v %v", input, output)
		}

		fmt.Println(command)

		cmd := exec.Command("/bin/bash", "-c", command)
		errStart := cmd.Start()
		if errStart != nil {
			return errStart
		}

		errWait := cmd.Wait()
		if errWait != nil {
			return errWait
		}
	}

	return nil
}

func RemoveShell(filePath string, isDir bool) error {
	_, errStat := os.Stat(filePath)
	if os.IsNotExist(errStat) {
		return nil
	}

	if runtime.GOOS == "windows" {
		command := ""
		if isDir {
			command = fmt.Sprintf("rd /S /Q %v", filePath)
		} else {
			command = fmt.Sprintf("del /F /S /Q %v", filePath)
		}

		fmt.Println(command)

		cmd := exec.Command("cmd", "/C", command)
		errStart := cmd.Start()
		if errStart != nil {
			return errStart
		}

		errWait := cmd.Wait()
		if errWait != nil {
			return errWait
		}
	} else {
		command := ""
		if isDir {
			command = fmt.Sprintf("rm -rf %v", filePath)
		} else {
			command = fmt.Sprintf("rm -f %v", filePath)
		}

		fmt.Println(command)

		cmd := exec.Command("/bin/bash", "-c", command)
		errStart := cmd.Start()
		if errStart != nil {
			return errStart
		}

		errWait := cmd.Wait()
		if errWait != nil {
			return errWait
		}
	}

	return nil
}

func ReadLineAt(file string, index int64) (string, error) {
	stat, errStat := os.Stat(file)
	if os.IsNotExist(errStat) {
		return "", nil
	}

	if errStat != nil {
		return "", errStat
	}

	if stat.Size() == 0 {
		return "", nil
	}

	f, errOpen := os.Open(file)
	if errOpen != nil {
		return "", errOpen
	}

	start := stat.Size()
	for {
		start = start - 1
		buf := make([]byte, 1)
		_, errRead := f.ReadAt(buf, start)
		if errRead == io.EOF {
			break
		}

		if errRead != nil {
			return "", errRead
		}

		if buf[0] == '\r' || buf[0] == '\n' {
			index = index - 1
			if index == 0 {
				break
			}
		}

		if start == 0 {
			break
		}
	}

	result := make([]byte, stat.Size()-start)
	_, errRead := f.ReadAt(result, start)
	if errRead != nil && errRead != io.EOF {
		return "", errRead
	}

	return string(result[:]), nil
}

func CreateDir(path string) error {
	_, errDest := os.Stat(path)
	if errDest == nil {
		return nil
	}
	if os.IsNotExist(errDest) {
		errDir := os.MkdirAll(path, os.ModeDir | os.ModePerm)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}

func IsExist(path string) (bool, error) {
	_, errDest := os.Stat(path)
	if errDest == nil {
		return true, nil
	}
	if os.IsNotExist(errDest) {
		return false, nil
	}
	return true, nil
}

func SaveCsv(fileName string, header []string, data [][]string) error {
	createFile, errCreateFile := os.Create(fileName)
	if errCreateFile != nil {
		return errCreateFile
	}

	defer createFile.Close()

	buf := new(bytes.Buffer)
	r2 := csv.NewWriter(buf)
	r2.Write(header)

	for _, d := range data {
		if d == nil {
			continue
		}

		r2.Write(d)
	}

	r2.Flush()
	createFile.WriteString(buf.String())

	return nil
}
