package main

type KeyValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Project struct {
	Name      string `json:"name"`
	Filename  string `json:"filename"`
	Workspace string `json:"workspace"`
	Target    string `json:"target"`
	Jboss     string `json:"jboss"`
	Java      string `json:"java"`
	Alias     string `json:"alias"`
	Flow      int    `json:"flow"`
}

type RawConfig struct {
	Workspaces []KeyValue `json:"workspaces"`
	Java       []KeyValue `json:"java"`
	Projects   []Project  `json:"projects"`
	Jboss      []KeyValue `json:"jboss"`
}

type Config struct {
	java       map[string]string
	jboss      map[string]string
	workspaces map[string]string
}

type Repository struct {
	projects map[string]Project
}

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
