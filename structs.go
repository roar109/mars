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
