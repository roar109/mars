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
	projects         = new(Projects)
	projectAlias     = flag.String("p", "", "Project Alias")
	skipDeploy       = flag.Bool("skipDeploy", false, "Skip deployment phase")
	skipArtifactCopy = flag.Bool("skipArtifactCopy", false, "Skip artifact copy to deployment folder")
	gitAction        = flag.String("git", "", "Execute git command")
	help             = flag.Bool("help", false, "Show available commands")
	configFileName   = flag.String("file", "config.json", "The path to the json configuration file, if not specify takes config.json by default")
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

	for _, proj := range *projects {
		if proj.Alias != "" {
			fmt.Printf("[%s] %s\n", proj.Alias, proj.Name)
		}
	}
	fmt.Println("***********************************")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nEnter alias: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(" ")

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

	flow := GetFlow(project.Flow)

	for _, stage := range states {
		if flow.Can(stage) {
			flow.Event(stage, project)
		}
	}
}

func printHelp() {
	fmt.Println("Help:")
	fmt.Println("-p=project-alias\tRun for a project")
	fmt.Println("-skipArtifactCopy\tDo not copy the generated artifact")
	fmt.Println("-skipDeploy\t\tIf enabled, do not start the java container")
	fmt.Println("-git=command\t\tIf available, run git commands before build")
	fmt.Println("-file=myfile.json\t\tThe path to the json configuration file, if not specify takes config.json by default")
}
