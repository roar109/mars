package main

import (
	"log"
	"os"
	"os/exec"
)

func executeGitAction(rootPath, action string) (status int, err error) {
	log.Println("Path to run Git on: " + rootPath)
	cmd := exec.Command("git", action)
	cmd.Env = os.Environ()
	/*
		for _, element := range cmd.Env {
			if strings.LastIndex(element, "GIT_EXEC_PATH") >= 0 {
				fmt.Println(element)
			}
		}
	*/
	return executeCommand(cmd)
}
