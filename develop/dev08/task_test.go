package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		out   string
	}{
		{
			name:  "cd",
			input: "cd ../dev03",
			out:   "",
		},
		{
			name:  "echo",
			input: "echo aaa",
			out:   "aaa\n",
		},
		{
			name:  "pwd",
			input: "pwd",
			out:   "/mnt/c/Users/Cybernet1c/Repository/WB_L2/WB_L2/develop/dev03\n",
		},
	}

	sh := Shell{}

	for _, tst := range tests {
		newBuf := bytes.NewBuffer(nil)
		sh.Out = newBuf

		err := sh.ProcessCommand(tst.input)
		if err != nil {
			t.Errorf("unexpected error %s\n", err)
		}

		assert.Equal(t, tst.out, newBuf.String(), "not equal")
	}
}

func TestProcessPipeline(t *testing.T) {
	tests := []struct {
		name  string
		input string
		out   string
	}{
		{
			name:  "pipeline",
			input: "go run ../dev03/task.go ../dev03/test.txt -k 3 -n -r -u | grep bud",
			out:   "123 bud f\r\n\n",
		},
	}

	sh := Shell{}

	for _, tst := range tests {
		newBuf := bytes.NewBuffer(nil)
		sh.Out = newBuf

		err := sh.ProcessPipeline(tst.input)
		if err != nil {
			t.Errorf("unexpected error %s\n", err)
		}

		assert.Equal(t, tst.out, newBuf.String(), "not equal")
	}
}
