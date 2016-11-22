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

	projectsLocal := make(Projects, len(rawConfig.Projects))

	//Parse projects array
	for _, proj := range rawConfig.Projects {
		repository.projects[strings.ToLower(proj.Alias)] = proj
		projectsLocal = append(projectsLocal, proj)
	}

	sort.Sort(projectsLocal)
	projects = &projectsLocal
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
