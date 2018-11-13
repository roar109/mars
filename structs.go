package main

//KeyValue representation of a json key:value
type KeyValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//Project Representation of a project
type Project struct {
	Name      string `json:"name"`
	Filename  string `json:"filename"`
	Workspace string `json:"workspace"`
	Target    string `json:"target"`
	Jboss     string `json:"jboss"`
	Java      string `json:"java"`
	Alias     string `json:"alias"`
	Maven     string `json:"maven"`
	Flow      int    `json:"flow"`
}

//RawConfig Representation of the given json file with the configuration
type RawConfig struct {
	Workspaces []KeyValue `json:"workspaces"`
	Java       []KeyValue `json:"java"`
	Projects   []Project  `json:"projects"`
	Jboss      []KeyValue `json:"jboss"`
	Maven      []KeyValue `json:"maven"`
	Sets       []Set      `json:"sets"`
}

//Config Cache or configurations with values parsed
type Config struct {
	java       map[string]string
	jboss      map[string]string
	workspaces map[string]string
	maven      map[string]string
}

//Set detail
type Set struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

//Repository Cache of the projects
type Repository struct {
	projects map[string]Project
}

//Projects Structure to compare projects and order them
type Projects []Project

func (projects Projects) Len() int {
	return len(projects)
}

func (projects Projects) Less(i, j int) bool {
	return projects[i].Alias < projects[j].Alias
}

func (projects Projects) Swap(i, j int) {
	projects[i], projects[j] = projects[j], projects[i]
}
