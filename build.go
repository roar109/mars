package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func build(project *Project) (status int, err error) {
	path := filepath.Join(config.workspaces[project.Workspace], project.Name, "pom.xml")
	mvnPath := filepath.Join(config.maven[project.Maven], "mvn")

	log.Println("Building...")
	cmd := exec.Command(mvnPath, "-f", path, "clean", "install")
	var env = os.Environ()

	env = addEnvVariable(env, "JAVA_HOME", config.java[project.Java])
	env = addEnvVariable(env, "JBOSS_HOME", config.jboss[project.Jboss])
	cmd.Env = env
	return executeCommand(cmd)
}
