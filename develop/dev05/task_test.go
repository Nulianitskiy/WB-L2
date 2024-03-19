package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrep(t *testing.T) {
	args := [][]string{
		{"-A", "2", "-B", "1", "str.w", "./tests/test.txt"},
		{"-C", "3", "3[dsa]", "./tests/test.txt"},
		{"-c", "-i", "sf", "./tests/test.txt"},
		{"-v", "-n", "qwe", "./tests/test.txt"},
		{"-n", "s*f", "./tests/test.txt"},
		{"-n", "-F", "s*f", "./tests/test.txt"},
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

		assert.Equal(t, expOutpur, curOutput, "files not equal")

	}

}
