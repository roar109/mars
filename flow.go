package main

import (
	"log"

	"github.com/looplab/fsm"
)

var states = []string{"build", "copy", "deploy"}

func GetFlow(flow int) *fsm.FSM {
	if flow == 1 {
		return NewBuildFlow()
	} else if flow == 2 {
		return NewBuildAndCopyFlow()
	} else if flow == 3 {
		return NewBuildAndDeployFlow()
	} else {
		return NewBuildAndDeployFlow()
	}
}

func NewBuildFlow() *fsm.FSM {
	log.Println("Using Build flow")
	return fsm.NewFSM(
		"inactive",
		fsm.Events{
			{Name: "build", Src: []string{"inactive"}, Dst: "inactive"},
		},
		fsm.Callbacks{
			"build": Build,
		},
	)
}

func NewBuildAndCopyFlow() *fsm.FSM {
	log.Println("Using Build and Copy flow")
	return fsm.NewFSM(
		"inactive",
		fsm.Events{
			{Name: "build", Src: []string{"inactive"}, Dst: "build"},
			{Name: "copy", Src: []string{"build"}, Dst: "inactive"},
		},
		fsm.Callbacks{
			"build": Build,
			"copy":  CopyArtifact,
		},
	)
}

func NewBuildAndDeployFlow() *fsm.FSM {
	log.Println("Using Build and Deploy flow")
	return fsm.NewFSM(
		"inactive",
		fsm.Events{
			{Name: "build", Src: []string{"inactive"}, Dst: "build"},
			{Name: "copy", Src: []string{"build"}, Dst: "copy"},
			{Name: "deploy", Src: []string{"copy"}, Dst: "inactive"},
		},
		fsm.Callbacks{
			"build":  Build,
			"copy":   CopyArtifact,
			"deploy": Deploy,
		},
	)
}

func Build(e *fsm.Event) {
	project := e.Args[0].(*Project)
	status, _ := build(project)
	if status > 0 {
		log.Fatal("Process failed")
	}
}

func CopyArtifact(e *fsm.Event) {
	if !*skipArtifactCopy {
		project := e.Args[0].(*Project)
		log.Println("Copying the file")
		copyArtifact(project)
	}
}

func Deploy(e *fsm.Event) {
	if !*skipArtifactCopy && !*skipDeploy {
		project := e.Args[0].(*Project)
		log.Println("Deploy not enabled")
		deploy(project)
		//jboss := filepath.Join(config.jboss[project.Jboss], "bin", "standalone")
	}
}
