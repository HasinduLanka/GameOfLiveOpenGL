package main

import (
	"bufio"
	"io"
	"os"
)

func DeleteFiles(name string) {
	os.RemoveAll(name)
}

func NewLineProvider_FromFile(filename string) (*LineProvider, error) {
	file, err := LoadFileToIOReader(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	LP := &LineProvider{make(chan string, 128)}

	go func() {
		for scanner.Scan() {
			LP.ch <- scanner.Text()
		}

		close(LP.ch)
		_ = file.Close()
	}()

	return LP, nil
}

func LoadFileToIOReader(filename string) (io.ReadCloser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Todo : optimize for bulk write
func AppendChan(filename string, ch chan []byte) {
	for line := range ch {
		AppendFile(filename, line)
	}
}

func AppendLine(filename string, line string) {
	AppendFile(filename, []byte(line+"\n"))
}

func AppendFile(filename string, content []byte) {
	F, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	CheckError(err)
	F.Write(content)
	F.Close()
}

func WriteFile(filename string, content []byte) {
	F, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	CheckError(err)
	F.Write(content)
	F.Close()
}
