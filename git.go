package main

import (
	"log"
	"os"
	"os/exec"
)

func executeGitAction(rootPath, action string) (status int, err error) {
	log.Println("Path to run git on: " + rootPath)
	cmd := exec.Command("git", "-C", rootPath, action)
	cmd.Env = os.Environ()
	return executeCommand(cmd)
}
