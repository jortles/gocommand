package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func commands(cmdStr string) error {
	cmdStr = strings.TrimSuffix(cmdStr, "\n")
	argsStr := strings.Fields(cmdStr)

	// Command Strings
	if len(argsStr) == 0 {
		fmt.Printf("\r")
		return nil
	}
	if argsStr[0] == "exit" {
		os.Exit(0)
	}

	// Used to run any cmd.exe/bash OS command
	if argsStr[0] == "c" {
		cmd := exec.Command(argsStr[1], argsStr[2:]...)
		cmd.Stdout = os.Stdout
		cmd.Run()
		return nil
	}
	if argsStr[0] == "cd" {
		cDir := os.Chdir(argsStr[1])
		if cDir == nil {
			return cDir
		}
		fmt.Print("Not a directory!\r\n")
	}

	// Windows only
	if argsStr[0] == "clear" {
		clr := exec.Command("cmd", "/c", "cls")
		clr.Stdout = os.Stdout
		clr.Run()
		return nil
	}

	// Example Commands
	if argsStr[0] == "pwd" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Print(wd)
		fmt.Printf("\r\n")
		return nil
	}
	if argsStr[0] == "id" {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Print(usr.Username)
		fmt.Printf("\r\n")
		return nil
	}
	if argsStr[0] == "ls" {
		list, err := ioutil.ReadDir(".")
		if err != nil {
			panic(err)
		}

		for _, file := range list {
			fmt.Println(file.Name())
		}
		fmt.Printf("\r")
		return nil
	}
	if argsStr[0] == "hostname" {
		host, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		fmt.Print(host)
		fmt.Printf("\r\n")
		return nil
	}
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("~$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		commands(cmdString)
	}
}
