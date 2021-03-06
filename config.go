package main

import (
	"os"
	"sort"
	"strings"
)

//Parse the raw json configuration
func parseProjectsConfig() {
	rawConfig := readConfigFile()

	//Fill maps
	config.maven = make(map[string]string, len(rawConfig.Maven))
	config.java = make(map[string]string, len(rawConfig.Java))
	config.jboss = make(map[string]string, len(rawConfig.Jboss))
	config.workspaces = make(map[string]string, len(rawConfig.Workspaces))
	repository.projects = make(map[string]Project, len(rawConfig.Projects))

	*sets = rawConfig.Sets

	for _, jav := range rawConfig.Java {
		config.java[jav.Name] = getSystemPropOrValue(jav.Value)
	}

	for _, jb := range rawConfig.Jboss {
		config.jboss[jb.Name] = getSystemPropOrValue(jb.Value)
	}

	for _, ws := range rawConfig.Workspaces {
		config.workspaces[ws.Name] = getSystemPropOrValue(ws.Value)
	}

	for _, mv := range rawConfig.Maven {
		config.maven[mv.Name] = getSystemPropOrValue(mv.Value)
	}

	projectsLocal := make(Projects, len(rawConfig.Projects))

	//Parse projects array
	for _, proj := range rawConfig.Projects {
		repository.projects[strings.ToLower(proj.Alias)] = proj
		projectsLocal = append(projectsLocal, proj)
	}

	sort.Sort(projectsLocal)
	projects = &projectsLocal
}

//Try to get a system variable with the given string, if not return the same string
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
