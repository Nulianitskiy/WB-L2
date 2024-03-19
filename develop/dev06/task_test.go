package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCut(t *testing.T) {
	args := [][]string{
		{"-f", "1,3", "-d", ":", "./tests/test.txt"},
		{"-f", "-2", "-d", ":", "./tests/test.txt"},
		{"-f", "1", "-d", ":", "-s", "./tests/test.txt"},
		{"-f", "1-3", "-d", " ", "-s", "./tests/test.txt"},
		{"-f", "2-", "-d", ":", "-s", "./tests/test.txt"},
	}

	for i, v := range args {
		old := os.Stdout
		filename := fmt.Sprintf("./tests/test%d.txt", i)
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Ошибка при создании временного файла:", err)
			return
		}

		os.Stdout = file

		err = Start(v)
		if err != nil {
			t.Errorf(err.Error())
		}

		f, err := os.Open(file.Name())
		if err != nil {
			t.Error(err.Error())
		}
		curOutput, err := io.ReadAll(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		expectedFilename := fmt.Sprintf("./tests/expected%d.txt", i+1)
		f2, err := os.Open(expectedFilename)
		if err != nil {
			t.Error(err.Error())
		}
		expOutpur, err := io.ReadAll(f2)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		os.Stdout = old

		assert.Equal(t, string(expOutpur), string(curOutput), "files not equal")

	}
}
