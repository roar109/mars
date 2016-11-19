package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func executeCommand(cmd *exec.Cmd) (code int, execErr error) {
	code = 0
	// Attach os  out to command out
	cmd.Stdout = os.Stdout

	printCommand(cmd)
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		printError(err)

		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
			code = waitStatus.ExitStatus()
			execErr = err
		}
	}
	return code, execErr
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func addEnvVariable(array []string, envirName string, envValue string) []string {
	var indexEln = -1

	for index, element := range array {
		if strings.LastIndex(element, envirName+"=") >= 0 {
			indexEln = index
		}
	}
	log.Println("Setting " + envirName + "=" + envValue)
	//If not found in existing env we need to add it
	if indexEln >= 0 {
		array[indexEln] = envirName + "=" + envValue
	} else {
		array = append(array, envirName+"="+envValue)
	}
	return array
}
