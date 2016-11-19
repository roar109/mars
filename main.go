package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	config           = new(Config)
	repository       = new(Repository)
	projectAlias     = flag.String("p", "", "Project Alias")
	skipDeploy       = flag.Bool("skipDeploy", false, "Skip deployment phase")
	skipArtifactCopy = flag.Bool("skipArtifactCopy", false, "Skip artifact copy to deployment folder")
	gitAction        = flag.String("git", "", "Execute git command")
	help             = flag.Bool("help", false, "Show available commands")
)

func init() {
	flag.Parse()
	parseProjectsConfig()
}

func main() {
	if val, ok := projectAliasExists(*projectAlias); ok {
		buildAndDeploy(&val)
		return
	}

	if *help {
		printHelp()
		return
	}

	fmt.Println("Output:\n[alias] Project Name")
	fmt.Println("\n********* Projects ****************")

	for k, v := range repository.projects {
		fmt.Printf("[%s] %s\n", k, v.Name)
	}
	fmt.Println("***********************************")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nEnter alias: ")
	text, _ := reader.ReadString('\n')

	if val, ok := projectAliasExists(text); ok {
		buildAndDeploy(&val)
	} else {
		fmt.Println("No valid option")
	}

}

func buildAndDeploy(project *Project) {
	//Pre steps
	if *gitAction != "" {
		path := filepath.Join(config.workspaces[project.Workspace], project.Name)
		status, _ := executeGitAction(path, *gitAction)
		if status > 0 {
			log.Fatal("Process failed")
		}
	}

	//Build artifact
	status, _ := build(project)
	if status > 0 {
		log.Fatal("Process failed")
	}

	//Copy the artifact
	if !*skipArtifactCopy {
		log.Println("Copying the file")
		copyArtifact(project)
	}

	//Run jboss with artifact
	if !*skipArtifactCopy && !*skipDeploy {
		log.Println("Deploy not enabled")
		deploy()
		//jboss := filepath.Join(config.jboss[project.Jboss], "bin", "standalone")
	}
}

func copyArtifact(project *Project) {
	artifact := filepath.Join(config.workspaces[project.Workspace], project.Name, project.Target, project.Filename)
	jbossDeploymentFolder := filepath.Join(config.jboss[project.Jboss], "standalone", "deployments", project.Filename)

	fmt.Printf("==> Copying %s to %s\n", artifact, jbossDeploymentFolder)

	err := CopyFile(artifact, jbossDeploymentFolder)
	if err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	fmt.Println("Help:")
	fmt.Println("-p=project-alias\tRun for a project")
	fmt.Println("-skipArtifactCopy\tDo not copy the generated artifact")
	fmt.Println("-skipDeploy\t\tIf enabled, do not start the java container")
	fmt.Println("-git=command\t\tIf available, run git commands before build")
}
