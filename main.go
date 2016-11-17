package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	config           = new(Config)
	repository       = new(Repository)
	projectAlias     = flag.String("p", "", "Project Alias")
	skipDeploy       = flag.Bool("skipDeploy", false, "Skip deployment phase")
	skipArtifactCopy = flag.Bool("skipArtifactCopy", false, "Skip artifact copy to deployment folder")
	gitAction        = flag.String("git", "", "Pull latest changes from git")
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

	if !*skipArtifactCopy {
		//Copy the artifact
		log.Println("Copying the file")
		copyArtifact(project)
	}

	//Run jboss with artifact
	if !*skipArtifactCopy && !*skipDeploy {
		log.Println("Deploy not enabled")
		//jboss := filepath.Join(config.jboss[project.Jboss], "bin", "standalone")
		//log.Println(jboss)
	}
}

func build(project *Project) (status int, err error) {
	path := filepath.Join(config.workspaces[project.Workspace], project.Name, "pom.xml")

	log.Println("Building...")
	cmd := exec.Command("mvn", "-f", path, "clean", "install")
	var env = os.Environ()

	//TODO User pointers/slices instead of copying the array
	env = addEnvVariable(env, "JAVA_HOME", config.java[project.Java])
	env = addEnvVariable(env, "JBOSS_HOME", config.jboss[project.Jboss])
	cmd.Env = env
	return executeCommand(cmd)
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

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func readConfigFile() *RawConfig {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var jsontype RawConfig
	err := json.Unmarshal(file, &jsontype)
	if err != nil {
		log.Fatal(err)
	}
	return &jsontype
}

//Parse the raw json configuration
func parseProjectsConfig() {
	rawConfig := readConfigFile()

	//Fill maps
	config.java = make(map[string]string, len(rawConfig.Java))
	config.jboss = make(map[string]string, len(rawConfig.Jboss))
	config.workspaces = make(map[string]string, len(rawConfig.Workspaces))
	repository.projects = make(map[string]Project, len(rawConfig.Projects))

	for _, jav := range rawConfig.Java {
		config.java[jav.Name] = getSystemPropOrValue(jav.Value)
	}

	for _, jb := range rawConfig.Jboss {
		config.jboss[jb.Name] = getSystemPropOrValue(jb.Value)
	}

	for _, ws := range rawConfig.Workspaces {
		config.workspaces[ws.Name] = getSystemPropOrValue(ws.Value)
	}

	//Parse projects array
	for _, proj := range rawConfig.Projects {
		repository.projects[strings.ToLower(proj.Alias)] = proj
	}
}

func getSystemPropOrValue(prop string) string {
	val := os.Getenv(prop)
	if val != "" {
		prop = val
	}
	return prop
}

func projectAliasExists(projAlias string) (Project, bool) {
	val, ok := repository.projects[strings.TrimSpace(projAlias)]
	return val, ok
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
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
