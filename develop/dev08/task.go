package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Shell struct {
	Out io.Writer
}

func ShellStart() {
	shell := Shell{
		Out: os.Stdout,
	}
	shell.Start()
}

func (sh *Shell) Start() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}

		text := scanner.Text()

		if strings.Contains(text, "|") {
			if err := sh.ProcessPipeline(text); err != nil {
				fmt.Fprintln(sh.Out, err.Error())
			}
		} else {
			err := sh.ProcessCommand(text)
			if err != nil {
				fmt.Fprintln(sh.Out, err.Error())
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(sh.Out, "Ошибка при считывании:", err)
	}
}

func (sh *Shell) ProcessPipeline(text string) error {
	cmds := strings.Split(text, " | ")
	execs := make([]*exec.Cmd, 0, len(cmds))
	for _, v := range cmds {
		args := strings.Split(v, " ")
		cmd := exec.Command(args[0], args[1:]...)
		execs = append(execs, cmd)
	}

	for i := range execs {
		if i > 0 {
			stdout, err := execs[i-1].Output()
			if err != nil {
				return err
			}
			b := bytes.NewReader(stdout)
			execs[i].Stdin = b
		}
		if i == len(execs)-1 {
			stdout, err := execs[i].Output()
			if err != nil {
				return err
			}
			//fmt.Println(string(stdout))
			fmt.Fprintln(sh.Out, string(stdout))
		}
	}

	return nil
}

func (sh *Shell) ProcessCommand(text string) error {
	cmds := strings.Split(text, " ")
	//fmt.Println(cmds)
	switch cmds[0] {
	case "cd":
		if len(cmds) > 1 {
			return os.Chdir(cmds[1])
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			return os.Chdir(home)
		}
	case "pwd":
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}

		fmt.Fprintln(sh.Out, currentDir)
	case "echo":
		for i := 1; i < len(cmds); i++ {
			fmt.Fprint(sh.Out, cmds[i])
			if i != len(cmds)-1 {
				fmt.Fprint(sh.Out, " ")
			}
		}
		fmt.Fprintln(sh.Out)
		return nil
	case "kill":
		if len(cmds) < 2 {
			return fmt.Errorf("not enough args")
		}
		pid := cmds[1]

		pidInt, err := strconv.Atoi(pid)
		if err != nil {
			return err
		}

		process, err := os.FindProcess(pidInt)
		if err != nil {
			return err
		}

		// Убиваем процесс
		err = process.Kill()

		return err
	case "ps":
		cmd := exec.Command("ps")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Ошибка выполнения команды:", err)
			return err
		}
		fmt.Fprintln(sh.Out, string(output))
	case "quit":
		fmt.Fprintln(sh.Out, "quitting")
		os.Exit(0)
	default:
		cmd := exec.Command(cmds[0], cmds[1:]...)

		// Выполнение команды
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Ошибка выполнения команды:", err)
			return err
		}

		// Вывод результата
		fmt.Fprintln(sh.Out, string(output))
	}

	return nil

}

func main() {
	ShellStart()
}
